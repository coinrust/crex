package crex

import (
	"fmt"
	"github.com/coinrust/crex/utils"
)

var idGen *utils.IdGenerate

func SetIdGenerate(g *utils.IdGenerate) {
	idGen = g
}

func GenOrderId() string {
	id := idGen.Next()
	return fmt.Sprint(id)
}
