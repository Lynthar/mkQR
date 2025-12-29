package cli

import (
	"fmt"

	"github.com/Lynthar/mkQR/internal/encoder"
	"github.com/spf13/cobra"
)

var (
	geoLat   float64
	geoLng   float64
	geoQuery string
)

var geoCmd = &cobra.Command{
	Use:   "geo",
	Short: "Generate QR code for geographic location",
	Long: `Generate a QR code containing geographic coordinates.

When scanned, this will open the location in a maps application.

Examples:
  mkqr geo --lat 40.7128 --lng -74.0060
  mkqr geo --lat 39.9042 --lng 116.4074 --query "Beijing"
  mkqr geo --lat 31.2304 --lng 121.4737 --query "Shanghai Tower"`,
	RunE: runGeo,
}

func init() {
	geoCmd.Flags().Float64Var(&geoLat, "lat", 0, "Latitude [required]")
	geoCmd.Flags().Float64Var(&geoLng, "lng", 0, "Longitude [required]")
	geoCmd.Flags().StringVar(&geoQuery, "query", "", "Location name/query")

	geoCmd.MarkFlagRequired("lat")
	geoCmd.MarkFlagRequired("lng")

	rootCmd.AddCommand(geoCmd)
}

func runGeo(cmd *cobra.Command, args []string) error {
	geo := &encoder.Geo{
		Latitude:  geoLat,
		Longitude: geoLng,
		Query:     geoQuery,
	}

	content := geo.Encode()

	if !quiet {
		if geoQuery != "" {
			fmt.Fprintf(cmd.ErrOrStderr(), "Location: %s (%.4f, %.4f)\n", geoQuery, geoLat, geoLng)
		} else {
			fmt.Fprintf(cmd.ErrOrStderr(), "Location: %.4f, %.4f\n", geoLat, geoLng)
		}
	}

	return generateQR(content)
}
