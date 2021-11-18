package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

const (
	NasaBaseUri = "https://api.nasa.gov/mars-photos/api/v1"
	ApiKey      = "DEMO_KEY"
)

type NasaResponse struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	ID        int64  `json:"id"`
	Sol       int64  `json:"sol"`
	Camera    Camera `json:"camera"`
	ImgSrc    string `json:"img_src"`
	EarthDate string `json:"earth_date"`
	Rover     Rover  `json:"rover"`
}

type Camera struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	RoverID  int64  `json:"rover_id"`
	FullName string `json:"full_name"`
}

type Rover struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	LandingDate string `json:"landing_date"`
	LaunchDate  string `json:"launch_date"`
	Status      string `json:"status"`
}

type getResponse map[string][]string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Obtains images from the Curiosity Mars Rover for the last 10 days",
	Run:   getImagesCmd,
}

// getImagesCmd will either return cached images or will attempt to retrieve new images from the NASA api
func getImagesCmd(cmd *cobra.Command, args []string) {
	// defer saving the cache such that it happens after commands execute
	defer putCache()
	c := viper.Get(RoverCache).(cache)

	r := make(getResponse)

	/*
		Count down 10 days from the current date and pull images from cache.
		If none are available, then query the Nasa API
	*/
	currentDate := time.Now()
	for i := -9; i <= 0; i++ {
		newDate := currentDate.Add(time.Duration(i) * 24 * time.Hour).Format("2006-01-02")
		if _, ok := c[rover][camera][newDate]; ok {
			for _, image := range c[rover][camera][newDate] {
				r.AddImage(newDate, image)
			}
		} else {
			resp, err := getNasaImages(rover, camera, newDate)
			if err != nil {
				logrus.Fatalln("error returned from getNasaImages:", err)
			}
			// write image urls to both the response & the cache
			for _, image := range resp.Photos {
				c[rover][camera][newDate] = append(c[rover][camera][newDate], image.ImgSrc)
				r.AddImage(newDate, image.ImgSrc)
			}
		}
	}

	b, err := json.Marshal(r)
	if err != nil {
		logrus.Fatalln("error occurred marshalling getImages response", err)
	}
	// Print the final response to stdout
	logrus.Println(string(b))
}

// getNasaImages will attempt to query the Nasa rover images API and pull imageUrls for the given rover, camera, and earthDate
func getNasaImages(rover, camera, earthDate string) (*NasaResponse, error) {
	r := &NasaResponse{}

	// Form the URL using the net/url package to ensure proper url formatting
	url := fmt.Sprintf("%s/rovers/%s/photos?earth_date=%s&camera=%s&api_key=%s",
		NasaBaseUri,
		rover,
		earthDate,
		camera,
		ApiKey,
	)

	resp, err := Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting data from nasa:", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http response body:", err)
	}

	if err := json.Unmarshal(b, r); err != nil {
		return nil, fmt.Errorf("error unmarshalling json to nasa response:", err)
	}

	return r, nil
}

// AddImage attempts to add an image url to the response, with a maximum of three images
func (r getResponse) AddImage(date, imageUrl string) {
	// Check to see if this date has already been used
	if _, ok := r[date]; ok {
		// if so, then ensure that no more than 3 images have been added for the date in our response
		if len(r[date]) < 3 {
			logrus.Debugf("adding image for %s: %s", date, imageUrl)
			r[date] = append(r[date], imageUrl)
		} else {
			logrus.Debugf("max images met for date key %s", date)
		}
	} else {
		r[date] = append(r[date], imageUrl)
	}
}
