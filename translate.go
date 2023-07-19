package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {

	///////////////////////////////////////
	// Init OpenAi
	///////////////////////////////////////
	openaiKey, ok := os.LookupEnv("OPENAI_KEY")
	if !ok {
		fmt.Println("Opsi. OPENAI_KEY export is not set. Please set it. E.g. 'export OPENAI_KEY=mysecretkey'. Generate the key here: https://platform.openai.com/account/api-keys")
		os.Exit(1)
	}

	if openaiKey == "" {
		fmt.Println("folder is required")
		os.Exit(1)
	}

	client := openai.NewClient(openaiKey)

	///////////////////////////////////////
	// Parse command line parameters
	///////////////////////////////////////
	var folder string

	flag.StringVar(&folder, "folder", "", "Please specify folder that contains markdown files (ending in .md) to translate")

	// Parse the flags
	flag.Parse()

	// Check if flags were set
	if folder == "" {
		fmt.Println("Opsi. Parameter 'folder' is required!")
		os.Exit(1)
	}

	///////////////////////////////////////
	// Get files
	///////////////////////////////////////
	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".de.md") {
			files = append(files, path) // append the file to the slice
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", folder, err)
	}

	////////////////////////////////////////////////////
	// Translate all files that were not translated yet
	////////////////////////////////////////////////////
	for _, file := range files {

		// check if file with ending de exists => if NOT translate
		var translatedFileName = replacePostfixOrReturnOriginal(file, ".md", ".de.md")

		if translatedFileName != file && !fileExists(translatedFileName) {
			fmt.Println("translating " + file + " -> " + translatedFileName)

			content := readFile(file)
			translationResult, _ := translate(*client, content)
			writeFile(translationResult, translatedFileName)

		} else {
			fmt.Println("Exists already as translated version. Not translating " + file)
		}

	}

}

func translate(client openai.Client, content string) (string, bool) {

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Can you translate the following hugo markdown file to German - do not change any markdown or html formatting. Also make sure not to change the header keys between the '---' markings. \n\n" + content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("OpenAI error: %v\n", err)
		return "", true
	}

	return resp.Choices[0].Message.Content, false

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func replacePostfixOrReturnOriginal(original string, oldPostfix string, newPostfix string) string {
	if strings.HasSuffix(original, oldPostfix) {
		return original[:len(original)-len(oldPostfix)] + newPostfix
	}
	return original
}

func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func writeFile(content string, outputFileName string) {
	err := os.WriteFile(outputFileName, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
