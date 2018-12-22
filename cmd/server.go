package cmd

import (
	"github.com/dimanevelev/aqua/config"
	"github.com/dimanevelev/aqua/handler"
	"github.com/dimanevelev/aqua/persistence"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var sConfig = config.ServerConst{}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a server and listens on a provided port",
	Long: `This command will start a server an expose two api endpoints:
	- /api/v1/file - Receives and stores POST requests with file information payload of the format {"Path":"/foo/bar","Name":"foo.bar","Size":123}.
	- /api/v1/stats - Receives GET requests and will return statistics of the received files. 
		Result example: {"code":200,"data":{"TotalFiles":2,"MaxFile":{"Size":3,"Path":"/foo/bar.abc"},"AvgFileSize":1.5,"Extensions":[".abc",".txt"],"TopExtension":".txt"}}
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting server")

		configuration, err := config.NewServerConfig(sConfig)

		if err != nil {
			log.Fatal("Configuration error:\n", err.Error())
		}
		router := Routes(configuration)

		log.Println("Serving application at Port " + configuration.ServerConst.Port)

		if configuration.Port == "443" {
			log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
			return
		}
		log.Fatal(http.ListenAndServe(":"+configuration.ServerConst.Port, router))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar( &sConfig.Port,"port", "8080", "The server port. For https use 443.")
	serverCmd.PersistentFlags().StringVar( &sConfig.MySqlConf.Port, "db-port", "3306", "The MySql DB port.")
	serverCmd.PersistentFlags().StringVar( &sConfig.MySqlConf.URL,"db-url", "127.0.0.1", "The MySql DB url.")
	serverCmd.PersistentFlags().StringVarP( &sConfig.MySqlConf.DBName,"db-name", "n","files", "The MySql DB name.")
	serverCmd.PersistentFlags().StringVarP( &sConfig.MySqlConf.User,"db-username", "u","root", "The MySql DB username.")
	serverCmd.PersistentFlags().StringVarP( &sConfig.MySqlConf.Password,"db-password", "p","password", "The MySql DB password.")
	//todo add validation
}

func Routes(config config.ServerConfig) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/file", handler.NewFileHandler(persistence.Client{Client: config.Database}).Routes())
		r.Mount("/stats", handler.NewStatsHandler(persistence.Client{Client: config.Database}).Routes())
	})

	return router
}
