package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/dimanevelev/aqua/config"
	"github.com/dimanevelev/aqua/model"
	"github.com/dimanevelev/aqua/parallel"
	"github.com/dimanevelev/aqua/utils"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tConfig = config.TraverserConfig{}

// traversCmd represents the travers command
var traversCmd = &cobra.Command{
	Use:   "travers",
	Short: "Traverses the file system and sends statistics to a server",
	Long: `Traverses the file system and sends statistics to a server. The statistics are file name, path, and size.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileHandlerFunc := createFileHandlerFunc()
		producerConsumer := parallel.NewBounedRunner(tConfig.Threads, false)
		errorsQueue := utils.NewErrorsQueue(tConfig.Threads)
		go func() {
			defer producerConsumer.Done()
			err := filepath.Walk(tConfig.Path, func(path string, f os.FileInfo, err error) error {
				if !f.IsDir() && f.Mode() != os.ModeSymlink {
					taskData := TaskData{path, f, tConfig.Url}
					task := fileHandlerFunc(taskData)
					_, _ = producerConsumer.AddTask(task, errorsQueue.AddError)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}()
		producerConsumer.Run()
		err := errorsQueue.GetError()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(traversCmd)
	traversCmd.PersistentFlags().StringVarP( &tConfig.Url,"url", "u", "http://localhost:8080/api/v1/file", "The target of the requests.")
	traversCmd.PersistentFlags().StringVarP( &tConfig.Path, "path", "p", ".","The path to start traversing.")
	traversCmd.PersistentFlags().IntVarP( &tConfig.Threads, "threads", "t", 3,"Number of threads that will send the requests")
}

type TaskData struct {
	FilePath string
	FileInfo os.FileInfo
	Url      string
}

type fileContext func(data TaskData) parallel.TaskFunc

func createFileHandlerFunc() fileContext {
	return func(data TaskData) parallel.TaskFunc {
		return func(threadId int) (err error) {
			log.Println("[Thread #", threadId, "]- Sending FileInfo of file:", data.FilePath)
			fileInfo := model.FileInfo{
				Name:    data.FileInfo.Name(),
				Size:    data.FileInfo.Size(),
				Mode:    data.FileInfo.Mode(),
				ModTime: data.FileInfo.ModTime(),
				IsDir:   data.FileInfo.IsDir(),
			}
			file := model.File{Path: data.FilePath, FileInfo:fileInfo}
			body, err := json.Marshal(file)
			req, err := http.NewRequest("POST", data.Url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Println("[Thread #", threadId, "- Received", resp.Status, "when sending info on path", data.FilePath)
			}
			return
		}
	}
}
