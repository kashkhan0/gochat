package main

import (
	// "flag"
	"fmt"
	"github.com/buger/jsonparser"
	// "github.com/gorilla/websocket"
	// "html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	// "net/http"
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var clog = log.Println

// func clog(v ...interface{}) {
// 	log.Println(v...)
// }

func system(cmd string, args []string) {

	var out []byte
	var err error

	// args := []string{""}
	if out, err = exec.Command(cmd, args...).Output(); err != nil {

		fmt.Fprintln(os.Stderr, "There was an error: ", err)
		log.Println("prob", cmd)
		return
	}
	// fmt.Println("Successful",cmd)
	sha := string(out)
	fmt.Println(cmd, strings.Join(args[:], " "))
	fmt.Println(sha)
}

func getkvf(key string) string {

	fn := "./" + key + ".txt"

	// fmt.Print("reading", fn)
	if _, err := os.Stat(fn); err == nil {
		// path/to/whatever exists

		dat, err := ioutil.ReadFile(fn)
		check(err)
		// fmt.Print(string(dat))
		return string(dat)
	}

	return "nil"
}

func getkvfile(key string) string {

	fn := "./" + key

	// fmt.Print("reading", fn)
	if _, err := os.Stat(fn); err == nil {
		// path/to/whatever exists

		dat, err := ioutil.ReadFile(fn)
		check(err)
		// fmt.Print(string(dat))
		return string(dat)
	}

	return "nil"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setkvf(key string, val string) {
	d1 := []byte(val)
	err := ioutil.WriteFile(key+".txt", d1, 0644)
	check(err)
}

func append2file(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := io.WriteString(f, text)
	if err != nil {
		fmt.Println(n, err)
		return
	}
	f.Close()
}

func log2file(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := io.WriteString(f, text)
	if err != nil {
		fmt.Println(n, err)
		return
	}
	f.Close()
}

// func readintosliceo(filename string) []string {
// 	content, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		//Do something
// 	}
// 	lines := strings.Split(string(content), "\n")
// 	return lines
// }

func readintoslice(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func readintostring(path string) string {
	file, err := os.Open(path)
	if err != nil {
		clog(err, path)
		return ""
	}
	defer file.Close()
	b := bufio.NewReader(file)
	s := ""
	for {
		s1, err := b.ReadString('\n')
		if err == io.EOF {
			break
		}

		s += s1
	}
	return s
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeslicetofile(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func getjson(js []byte, key string) string {
	out := ""
	if value, err := jsonparser.GetString(js, key); err == nil {
		out = value
	}

	return out
}

func listfiles(dir string) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fi, err := os.Stat(f.Name())
		if err != nil {
			// Could not obtain stat, handle error
			return
		}
		fmt.Println(f.Name(), fi.Size())
	}
}

func thuman() string {

	t := time.Now()
	// p(t.Format(time.RFC3339))
	yy := fmt.Sprintf("%04d", t.Year())
	mo := fmt.Sprintf("%02d", t.Month())
	dd := fmt.Sprintf("%02d", t.Day())
	hh := fmt.Sprintf("%02d", t.Hour())
	mi := fmt.Sprintf("%02d", t.Minute())
	ss := fmt.Sprintf("%02d", t.Second())
	ns := fmt.Sprintf("%09d", t.Nanosecond())
	a := yy + mo + dd + hh + mi + ss + "." + ns
	// + s(t.Hour()) + s( t.Minute()) + s( t.Second())
	return a

}

func mkdir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		// Could not obtain stat, handle error
		system("mkdir", []string{"-p", dir})
	}
}

func atoi(a string) int {
	out, _ := strconv.Atoi(a)
	return out
}
func atof(a string) float64 {
	out, _ := strconv.ParseFloat(a, 64)
	return out
}

func randseq(n int) string {
	letterRunes := []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//
//
//========================================================================
//
//
// wslib
//
//
//========================================================================
//
//
