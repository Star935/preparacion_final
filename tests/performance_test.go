package tests

import (
    "context"
    "encoding/json"
    "math/rand"
    "net/http/httptest"
    "testing"
    "time"

    vegeta "github.com/tsenart/vegeta/v12/lib"
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "backend/handlers"
    "backend/models"
)

func randomString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}

func TestPerformanceCreateUsers(t *testing.T) {
    rand.Seed(time.Now().UnixNano())

    // === SETUP (timeout corto) ===
    ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelSetup()

    client, err := mongo.Connect(ctxSetup, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatal(err)
    }
    userColl := client.Database("perf_test_db").Collection("users")
    if err := userColl.Drop(ctxSetup); err != nil {
        t.Fatal(err)
    }

    // === LEVANTAR SERVER en httptest ===
    h := handlers.NewHandler(nil, userColl, nil)
    e := echo.New()
    e.POST("/users", h.CreateUser)
    ts := httptest.NewServer(e)
    defer ts.Close()

    // === ATTACK de Vegeta (10 s) ===
    rate := vegeta.Rate{Freq: 10, Per: time.Second}
    duration := 10 * time.Second
    attacker := vegeta.NewAttacker()
    var metrics vegeta.Metrics

    targeter := func() vegeta.Targeter {
        return func(tgt *vegeta.Target) error {
            u := models.User{
                Name:  randomString(12),
                Email: randomString(8) + "@test.com",
            }
            body, _ := json.Marshal(u)
            *tgt = vegeta.Target{
                Method: "POST",
                URL:    ts.URL + "/users",
                Body:   body,
                Header: map[string][]string{"Content-Type": {"application/json"}},
            }
            return nil
        }
    }()

    for res := range attacker.Attack(targeter, rate, duration, "perf-create-users") {
        metrics.Add(res)
    }
    metrics.Close()

    // === VERIFICACIONES ===
    if metrics.Success < 0.9 {
        t.Fatalf("esperaba ≥90%% de éxitos, pero fue %.2f%%", metrics.Success*100)
    }

    // crea un nuevo contexto para consulta (o usa context.Background())
    ctxQuery, cancelQuery := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelQuery()

    count, err := userColl.CountDocuments(ctxQuery, bson.M{})
    if err != nil {
        t.Fatal(err)
    }
    if count == 0 {
        t.Fatal("no se creó ningún usuario en la base de datos")
    }
    t.Logf("Se crearon %d usuarios con éxito", count)
}
