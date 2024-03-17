package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/realPointer/thumbnails-downloader-cli/pkg/thumbnail_v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	asyncDownload bool
	outputDir     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "thumbnails-downloader",
	Short: "Download YouTube video thumbnails",
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("Please provide at least one YouTube video link")
		}

		client, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer client.Close()

		thumbnailClient := thumbnail_v1.NewThumbnailServiceClient(client)

		if asyncDownload {
			var wg sync.WaitGroup
			for _, link := range args {
				wg.Add(1)
				go func() {
					defer wg.Done()
					downloadThumbnail(thumbnailClient, link, outputDir)
				}()
			}
			wg.Wait()
		} else {
			for _, link := range args {
				downloadThumbnail(thumbnailClient, link, outputDir)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&asyncDownload, "async", false, "Enable asynchronous downloading")
	rootCmd.Flags().StringVar(&outputDir, "output", ".", "Output directory for downloaded thumbnails")
}

func downloadThumbnail(client thumbnail_v1.ThumbnailServiceClient, link, outputDir string) {
	videoID := extractVideoID(link)
	if videoID == "" {
		fmt.Printf("Invalid YouTube link: %s\n", link)

		return
	}

	req := &thumbnail_v1.DownloadThumbnailRequest{
		VideoUrl: videoID,
	}

	resp, err := client.DownloadThumbnail(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed to download thumbnail for video %s: %v\n", videoID, err)

		return
	}

	filename := fmt.Sprintf("%s.jpg", videoID)
	outputPath := filepath.Join(outputDir, filename)

	err = os.WriteFile(outputPath, resp.ThumbnailData, 0o644)
	if err != nil {
		fmt.Printf("Failed to save thumbnail for video %s: %v\n", videoID, err)

		return
	}

	fmt.Printf("Downloaded thumbnail for video %s\n", videoID)
}

func extractVideoID(link string) string {
	// https://www.youtube.com/watch?v=-I26usMiuS4 -> -I26usMiuS4
	parts := strings.Split(link, "?v=")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}
