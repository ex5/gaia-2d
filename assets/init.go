package assets

import (
	"github.com/EngoEngine/engo/common"

	//"fmt"
	"encoding/json"
	//"errors"
	"io/ioutil"
	"math/rand"
	"os"
)

type Matter struct {
	ID int          `json:"id"`
	Type string     `json:"type"`
}

type Object struct {
	ID int          `json:"id"`
	Name string     `json:"name"`
	MatterID int    `json:"matter_id"`
	Amount float32  `json:"amount"`

	Matter *Matter
}

type Objects struct {
    Objects []*Object `json:"objects"`
}

type KindsOfMatter struct {
    KindsOfMatter []*Matter `json:"matter"`
}

var (
	LineHeight = 20
	FontSize = 16
	SpriteWidth = 32
	SpriteHeight = 32
	PreloadList = []string{
		"textures/chick_32x32.png",
		"tilemap/terrain-v7.png",
	}
	FullSpriteSheet *common.Spritesheet
	objects *Objects
	matter *KindsOfMatter
)

func InitAssets() {
	// Load the spritesheet
	FullSpriteSheet = common.NewSpritesheetFromFile("tilemap/terrain-v7.png", 32, 32)

	// Load objects
	objectsJsonFile, _ := os.Open("assets/meta/objects.json")
	defer objectsJsonFile.Close()
	byteValue, _ := ioutil.ReadAll(objectsJsonFile)
	json.Unmarshal(byteValue, &objects)

	// Load other related metadata
	metaJsonFile, _ := os.Open("assets/meta/matter.json")
	defer metaJsonFile.Close()
	byteValue, _ = ioutil.ReadAll(metaJsonFile)
	json.Unmarshal(byteValue, &matter)
	matterById := make(map[int]*Matter)
	for _, m := range matter.KindsOfMatter {
		matterById[m.ID] = m
	}
	// Fill in references
	for _, v := range objects.Objects {
		v.Matter = matterById[v.MatterID]
	}
}

func GetObjectsByType(matterType string) []*Object {
	var result []*Object
	for _, v := range objects.Objects {
		if v.Matter.Type == matterType {
			result = append(result, v)
		}
	}
	return result
}

func GetRandomObjectOfType(matterType string) *Object {
	objects := GetObjectsByType(matterType)
	return objects[rand.Intn(len(objects))]
}
