package seeder

import (
	"context"
	"flag"
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/bootstrap"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/kholidss/movie-fest-skilltest/internal/seeder"
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
	"log"
	"os"
)

var (
	flags = flag.NewFlagSet("db:seeder", flag.ExitOnError)

	help = flags.Bool("help", false, "print help")
	run  = flags.String("run", "", "run seeder process")
)

func DoSeeder() {
	var (
		ctx = context.Background()
	)

	flags.Usage = usage
	_ = flags.Parse(os.Args[2:])

	args := flags.Args()

	if (len(args) == 0 && *run == "") || *help {
		flags.Usage()
		return
	}

	cfg, err := config.LoadAllConfigs()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load configuration file: %v", err))
	}

	flags.Usage = usage

	db := bootstrap.RegistryMySQLDatabase(cfg)
	defer db.Close()

	//define repository
	repoUser := repositories.NewUserRepository(db)
	repoGenre := repositories.NewGenreRepository(db)

	//define seeder controller
	cs := seeder.NewSeedRun(cfg, repoUser, repoGenre)

	if *run == "admin" {
		cs.AdminData(ctx)
		return
	}
	if *run == "genres" {
		cs.GenresData(ctx)
		return
	}
	log.Fatal("command seeder not available")
}

func usage() {
	fmt.Println(usageCommands)
}

var (
	usageCommands = `
Commands:
    run                Run the seeder process following with seeder command
 	help               Show all available commands
`
)
