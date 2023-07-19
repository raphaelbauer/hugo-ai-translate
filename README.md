# Hugo AI Translate

Hugo (https://gohugo.io/) is a content management framework. Content is stored as markdown files (md). Sometimes you want to semi-automatically translate content into another language. For me this was translating English content to German. Wouldn't it be cool if this'd happen automagically? 

**NOTE: No warranty. Use at your own risk. Tokens will be sent to OpenAI and it will be billed against your OpenAI account.**

## What is this?
- A command line program that translates English Hugo Markdown files into German.
- If a myfile.md exists it will create a myfile.de.md file
- If a myfile.de.md exists it will not translate the file and assumes that has been done already.
- **Rule of thumb: With ChatGpt3.5 one post will cost around 2 cents to translate. It heavily depends on how long your content is of course.**

## Usage:
- Set the openai key via ```export OPENAI_KEY=mysecretkey```
- Run the app  ```go run translate.go -f /Users/me/my-hugo-project/content/```