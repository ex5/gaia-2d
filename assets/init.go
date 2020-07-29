package assets

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"bytes"
	"encoding/json"
	"fmt"
	"gogame/data"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

var (
	SpriteWidth  = 32
	SpriteHeight = 32
	// UI
	LineHeight              = 20
	FontURL                 = "fonts/arcade_n.ttf"
	FontSize        float64 = 14
	HoverInfoHeight float32 = 200
	HUDLayer        float32 = 1001
	HUDMarginL      float32 = 20
	HUDMarginT      float32 = 20

	creatures    *data.Creatures
	objects      *data.Objects
	resources    *data.Resources
	spritesheets *data.Spritesheets

	CreatureById    map[int]*data.Creature
	ObjectById      map[int]*data.Object
	ResourceById    map[int]*data.Resource
	ResourceByType  map[string]*data.Resource
	SpritesheetById map[int]*data.Spritesheet

	WorkDir string
)

func readJSON(path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	return byteValue
}

func loadSpritesheets() {
	log.Println(spritesheets)

	spritesheets.Loaded = make(map[int]*common.Spritesheet)
	for _, sheet := range spritesheets.Spritesheets {
		if err := engo.Files.Load(sheet.FilePath); err != nil {
			panic(err)
		}
		spritesheets.Loaded[sheet.ID] = common.NewSpritesheetFromFile(sheet.FilePath, SpriteWidth, SpriteHeight)
	}
}

func InitAssets() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	WorkDir = path

	engo.Files.SetRoot(WorkDir + "/assets")
	log.Println(engo.Files.GetRoot())

	// Load the font
	engo.Files.LoadReaderData(FontURL, bytes.NewReader(gosmallcaps.TTF))

	// Load the spritesheets
	byteValue := readJSON("assets/meta/spritesheets.json")
	json.Unmarshal(byteValue, &spritesheets)
	loadSpritesheets()

	// Load objects
	byteValue = readJSON("assets/meta/objects.json")
	json.Unmarshal(byteValue, &objects)

	// Load other related metadata
	byteValue = readJSON("assets/meta/resources.json")
	json.Unmarshal(byteValue, &resources)

	// Load creatures
	byteValue = readJSON("assets/meta/creatures.json")
	json.Unmarshal(byteValue, &creatures)

	// Prepare hashes for ease of access to loaded assets
	SpritesheetById = make(map[int]*data.Spritesheet)
	for _, s := range spritesheets.Spritesheets {
		SpritesheetById[s.ID] = s
	}
	ResourceById = make(map[int]*data.Resource)
	ResourceByType = make(map[string]*data.Resource)
	for _, r := range resources.Resources {
		ResourceById[r.ID] = r
		ResourceByType[r.Type] = r
	}
	ObjectById = make(map[int]*data.Object)
	for _, o := range objects.Objects {
		ObjectById[o.ID] = o
		o.Spritesheet = spritesheets.Loaded[o.SpritesheetID]
		o.Animations = SpritesheetById[o.SpritesheetID].Animations
	}
	CreatureById = make(map[int]*data.Creature)
	for _, c := range creatures.Creatures {
		CreatureById[c.ID] = c
	}
}

func GetResourceByID(resourceID int) *data.Resource {
	resource, ok := ResourceById[resourceID]
	if ok {
		return resource
	}
	panic(fmt.Sprintf("Resource '%d' could not be found", resourceID))
}

func GetResourceByType(resourceType string) *data.Resource {
	resource, ok := ResourceByType[resourceType]
	if ok {
		return resource
	}
	panic(fmt.Sprintf("Resource '%s' could not be found", resourceType))
}

func GetObjectById(objectID int) *data.Object {
	object, ok := ObjectById[objectID]
	if ok {
		return object
	}
	panic(fmt.Sprintf("Object '%d' could not be found", objectID))
}

func GetObjectsByType(resourceType string) []*data.Object {
	var result []*data.Object
	for _, v := range objects.Objects {
		if GetResourceByID(v.ResourceID).Type == resourceType {
			result = append(result, v)
		}
	}
	return result
}

func GetRandomObjectOfType(resourceType string) *data.Object {
	objects := GetObjectsByType(resourceType)
	return objects[rand.Intn(len(objects))]
}

func GetSpritesheetById(spritesheetID int) *common.Spritesheet {
	sheet, ok := spritesheets.Loaded[spritesheetID]
	if ok {
		return sheet
	}
	panic(fmt.Sprintf("Spritesheet '%d' does not appear to be loaded", spritesheetID))
}

func GetCreatureById(creatureID int) *data.Creature {
	creature, ok := CreatureById[creatureID]
	if ok {
		return creature
	}
	panic(fmt.Sprintf("Unknown creature '%d'", creatureID))
}
