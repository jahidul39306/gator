package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/google/uuid"
	"github.com/jahidul39306/gator/internal/config"
	"github.com/jahidul39306/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.commandNames[cmd.name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.commandNames[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("username is required")
	}
	username := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user '%s' does not exist", username)
	}
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User '%s' has been set\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("username is required")
	}
	username := cmd.arguments[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	cursor, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	s.cfg.SetUser(cursor.Name)
	fmt.Printf("User '%s' has been created\n", cursor.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete all users: %w", err)
	}
	fmt.Println("All users have been deleted")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get all users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching feed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i := range rssFeed.Channel.Items {
		rssFeed.Channel.Items[i].Title = html.UnescapeString(rssFeed.Channel.Items[i].Title)
		rssFeed.Channel.Items[i].Description = html.UnescapeString(rssFeed.Channel.Items[i].Description)
	}

	return &rssFeed, nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rssFeed, err := fetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	data, err := json.MarshalIndent(rssFeed, "", "  ")
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("Need two arguments")
	}
	name, url := cmd.arguments[0], cmd.arguments[1]
	username := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Fetching user: %w", err)
	}
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	cursor, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Creating feed: %w", err)
	}
	fmt.Printf("%v\n", cursor)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	allFeeds, err := s.db.GetAllFeedsWithUser(context.Background())
	if err != nil {
		return fmt.Errorf("fetching all feeds: %w", err)
	}
	fmt.Printf("%v\n", allFeeds)
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("url is required")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Fetching feed by url: %w", err)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Fetching user data: %w", err)
	}

	feedFollowsArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feed_follows, err := s.db.CreateFeedFollow(context.Background(), feedFollowsArgs)
	if err != nil {
		return fmt.Errorf("Creating feed follows: %w", err)
	}
	fmt.Printf("%v\n", feed_follows)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}
	defer db.Close()
	dbQueries := database.New(db)

	sta := state{cfg: cfg, db: dbQueries}

	cmds := commands{commandNames: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal(fmt.Errorf("not enough arguments were provided"))
	}
	cmd := command{name: args[0], arguments: args[1:]}
	err = cmds.run(&sta, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
