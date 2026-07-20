package stats

import(
	"spotify-stats/internal/spotify-data" 
)

type DistanceType int
const(
	ArtistPopularity DistanceType = iota
	SongPopularity
	Genre
	ImageColour
)

func CalculationStategy(calculationType DistanceType){
	switch(calculationType){
		case DistanceType.ArtistPopularity:
		case DistanceType.SongPopularity:
		case DistanceType.Genre:
		case DistanceType.ImageColour:
	}
}

// Data extractions
func ExtractAlbumData(){
	service := spotifydata.NewService()
	topOptions := spotifydata.TopOptions{
		spotifydata.TimeRange.ShortTerm,
		"50"
	}
	topTracks := service.FetchTopArtists()
	
	for _, track := range topTracks{
		
	}
}
func ExtractSongData(){}
func ExtractArtistData(){}
