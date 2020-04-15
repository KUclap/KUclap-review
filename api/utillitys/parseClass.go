package utillitys

import (
	"encoding/json"
	"fmt"
	_ "strings"
	"io/ioutil"
	"github.com/marsDev31/kuclap-backend/api/models"
)


func main(){
	var classes []models.Classes

	// Read json file
	data, err := ioutil.ReadFile("../classes.json")
    if err != nil {
      fmt.Print(err)
	}
	
	err = json.Unmarshal(data, &classes)
	if err != nil {
        fmt.Println("error:", err)
	}
	
	for _, class := range classes {
		fmt.Println("\"classID:\": " + "\"" + class.Value + "\",")
		// fmt.Println("\"nameTH:\": " + "\"" + strings.Fields(class.Lable) + "\"")
		// fmt.Println("\"nameEN:\": " + "\"" + class.Lable + "\"")
		fmt.Println("\"label:\": " + "\"" + class.Lable + "\"")
	}
}