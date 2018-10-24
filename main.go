package main
import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_"github.com/gin-gonic/gin/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"strconv"

	"io"
	"log"
	"net/http"
	"os"
)

var (
	USER=""
	PASS=""
	NAME=""
	ADDR=""
)

type(
	Player struct {
		Pid		  int		`json:"player_id"`
		Username  string    `json:"player_name"`
	}
	Game struct {
		Gameid	  int       `json:"game_id"`
		Playerid    int       `json:"player_id"`
		Score     int       `json:"score"`
		Date      int       `json:"date"`
		Kills     int       `json:"kills"`
	}
)

func main() {
	f, _ := os.Create("logs/gamejam.log")
	gin.DefaultWriter = io.MultiWriter(f)

	err := godotenv.Load()
	USER = os.Getenv("DB_USER")
	PASS = os.Getenv("DB_PASS")
	NAME = os.Getenv("GJAM_DB")
	ADDR = os.Getenv("DB_ADDR")

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET   ("/test", testFunc)
		api.GET	  ("/players", getPlayers)
		api.GET   ("/players/:id", getPlayer)
	}

	r.Run()
}

func testFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getPlayers(c *gin.Context) {
	limit := c.DefaultQuery("limit", "25")

	var query = "SELECT * FROM " + NAME + ".player LIMIT " + limit
	var players[] Player

	db, err := sql.Open("mysql", USER + ":" + PASS + "@tcp(" + ADDR + ")/" + NAME)
	if err != nil {
		fmt.Printf("Encountered error connecting to DB: " + err.Error() + "\n")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Connection to DB failed: " + err.Error() + "\n")
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("DB call failed: " + err.Error() + "\n")
	}

	for rows.Next() {
		var player Player
		err = rows.Scan(&player.Pid, &player.Username)
		players = append(players, player)
		fmt.Printf("Found player: " + strconv.Itoa(player.Pid) + " with name: " + player.Username + "\n")
	}

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H {
		"result": players,
		"count" : len(players),
	})

}

func getPlayer(c *gin.Context) {
	id := c.Param("id")
	var query = "SELECT * FROM " + NAME + ".player WHERE pid = ?"

	db, err := sql.Open("mysql", USER + ":" + PASS + "@tcp(" + ADDR + ")/" + NAME)
	if err != nil {
		fmt.Printf("Encountered error connecting to DB: " + err.Error() + "\n")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Connection to DB failed: " + err.Error() + "\n")
	}

	var player Player

	row := db.QueryRow(query, id)
	err = row.Scan(&player.Pid, &player.Username)

	if err != nil {
		fmt.Printf( "Encountered error finding player: " + err.Error() + "\n")
		c.JSON(http.StatusNotFound, gin.H {
			"result": "No players could be found.",
			"status_code": http.StatusNotFound,
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"result" : player,
			"count"  : 1,
		})
	}
}