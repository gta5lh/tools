package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/juju/errors"
)

const (
	serverID = "blue"
	maxLen   = 4
)

// 36 characters
var alphabet = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())

	count, err := generateStaticIDs(alphabet)
	if err != nil {
		log.Fatal(errors.Annotate(err, "generating static IDs failed"))
	}

	fmt.Println(fmt.Sprintf("%v IDs generated in %v", count, time.Since(start)))
}

func generateStaticIDs(alphabet []string) (int, error) {
	ids := generate(alphabet, maxLen)

	f, err := os.Create("generated_static_ids.csv")
	if err != nil {
		return 0, errors.Annotate(err, "creating file failed")
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(errors.Annotate(err, "closing file failed"))
		}
	}()

	w := bufio.NewWriter(f)
	if _, err := w.WriteString("static_id;server_id\n"); err != nil {
		return 0, errors.Annotate(err, "writing header failed")
	}

	for _, n := range ids {
		if _, err := w.WriteString(fmt.Sprintf("%s;%s\n", n, serverID)); err != nil {
			return 0, errors.Annotate(err, "writing static ID failed")
		}
	}

	if err := w.Flush(); err != nil {
		return 0, errors.Annotate(err, "flushing failed")
	}

	return len(ids), nil
}

func generate(alphabet []string, maxLen int) (res []string) {
	var ml []string
	ml = alphabet
	for _, symbol := range ml {
		res = append(res, symbol)
	}

	for z := 0; z < maxLen-1; z++ {
		var tmp []string

		for i := 0; i < len(alphabet); i++ {
			for k := 0; k < len(ml); k++ {
				tmp = append(tmp, ml[k]+alphabet[i])
			}
		}

		ml = tmp

		rand.Shuffle(len(tmp), func(i, j int) {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		})
		res = append(res, tmp...)
	}

	return res
}
