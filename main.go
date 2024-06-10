/*
**
** hamon
** Uses djb2 (1991) hash to map IP addresses to readable words.
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package main

import (
  "bufio"
  "flag"
  "fmt"
  "math/rand"
  "os"
  "strconv"
  "strings"
  "time"
  "unicode"
)


func init() {
  rand.Seed(time.Now().UnixNano())
}

func main() {
  forward := flag.Bool("forward", false, "Convert words to IP address")
  reverse := flag.Bool("reverse", false, "Convert IP address to words")
  full := flag.Bool("full", false, "Display all word options for each block")
  flag.Parse()

  if *forward && *reverse {
    fmt.Println("Please specify only one of -forward or -reverse.")
    return
  }

  if *forward {
    handleForward(flag.Args())
  } else if *reverse {
    handleReverse(flag.Args(), *full)
  } else {
    fmt.Println("Usage: hamon -forward \"word:word:word:word:word:word:word:word\" or -reverse \"ip_address\"")
  }
}

func handleForward(args []string) {
  if len(args) != 1 {
    fmt.Println("Usage: hamon -forward \"word:word:word:word:word:word:word:word\" or \"word.word.word.word\"")
    return
  }
  input := strings.ToLower(args[0])
  if strings.Contains(input, ":") {
    words := strings.Split(input, ":")
    generateIP(words, true)
  } else if strings.Contains(input, ".") {
    words := strings.Split(input, ".")
    generateIP(words, false)
  } else {
    fmt.Println("Error: Invalid format. Use colons for IPv6 or periods for IPv4.")
  }
}

func handleReverse(args []string, full bool) {
  if len(args) != 1 {
    fmt.Println("Usage: hamon -reverse \"ipv6_address\" or \"ipv4_address\"")
    return
  }
  ipAddress := strings.ToLower(args[0])
  if strings.Contains(ipAddress, ":") {
    reverseIP(ipAddress, true, full)
  } else if strings.Contains(ipAddress, ".") {
    reverseIP(ipAddress, false, full)
  } else {
    fmt.Println("Error: Invalid format. Use colons for IPv6 or periods for IPv4.")
  }
}

func djb2Hash(word string) uint32 {
  var hash uint32 = 5381
  for _, char := range word {
    hash = ((hash << 5) + hash) + uint32(char)
  }
  return hash
}

func getIPBlock(word string, isIPv6 bool) string {
  hash := djb2Hash(word)
  if isIPv6 {
    return fmt.Sprintf("%04x", hash%65536)
  }
  return fmt.Sprintf("%d", hash%256)
}

func loadWordMappings(filename string, isIPv6 bool) (map[string][]string, []string, error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, nil, err
  }
  defer file.Close()

  blockToWords := make(map[string][]string)
  var allWords []string
  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    word := strings.ToLower(scanner.Text())
    block := getIPBlock(word, isIPv6)
    blockToWords[block] = append(blockToWords[block], word)
    allWords = append(allWords, word)
  }

  if err := scanner.Err(); err != nil {
    return nil, nil, err
  }

  return blockToWords, allWords, nil
}

func generateIP(words []string, isIPv6 bool) {
  expectedLength := 8
  separator := ":"
  if !isIPv6 {
    expectedLength = 4
    separator = "."
  }

  if len(words) != expectedLength {
    fmt.Printf("Error: Please provide exactly %d words separated by %s.\n", expectedLength, separator)
    return
  }

  ipBlocks := make([]string, expectedLength)
  for i, word := range words {
    ipBlocks[i] = getIPBlock(word, isIPv6)
  }

  ipAddress := strings.Join(ipBlocks, separator)
  fmt.Println(ipAddress)
}

func reverseIP(ipAddress string, isIPv6 bool, full bool) {
  blocks := strings.Split(ipAddress, ":")
  expectedLength := 8
  if !isIPv6 {
    blocks = strings.Split(ipAddress, ".")
    expectedLength = 4
  }

  if len(blocks) != expectedLength {
    fmt.Printf("Error: Please provide a valid %s address with exactly %d blocks.\n", ipType(isIPv6), expectedLength)
    return
  }

  for _, block := range blocks {
    if !isValidBlock(block, isIPv6) {
      fmt.Printf("Error: Block %s is not a valid %s value.\n", block, blockType(isIPv6))
      return
    }
  }

  blockToWords, allWords, err := loadWordMappings("english-words/words_alpha.txt", isIPv6)
  if err != nil {
    fmt.Println("Error loading word mappings:", err)
    return
  }

  if full {
    printFull(blocks, blockToWords, allWords, isIPv6)
  } else {
    printRandom(blocks, blockToWords, allWords, isIPv6)
  }
}

func printFull(blocks []string, blockToWords map[string][]string, allWords []string, isIPv6 bool) {
  for i, block := range blocks {
    words, found := blockToWords[block]
    if !found {
      fmt.Printf("No words found for block %s. Generating fallback...\n", block)
      fallbackWord := generateFallbackWord(allWords, block, isIPv6)
      fmt.Printf("Words for block %d (%s): %v (fallback)\n", i+1, block, []string{fallbackWord})
    } else {
      fmt.Printf("Words for block %d (%s): %v\n", i+1, block, words)
    }
  }
}

func printRandom(blocks []string, blockToWords map[string][]string, allWords []string, isIPv6 bool) {
  var result []string
  for _, block := range blocks {
    words, found := blockToWords[block]
    if !found {
      fallbackWord := generateFallbackWord(allWords, block, isIPv6)
      result = append(result, fallbackWord)
    } else {
      randomWord := words[rand.Intn(len(words))]
      result = append(result, randomWord)
    }
  }

  separator := ":"
  if !isIPv6 {
    separator = "."
  }
  fmt.Println(strings.Join(result, separator))
}

func generateFallbackWord(allWords []string, block string, isIPv6 bool) string {
  for i := 0; ; i++ {
    for _, word := range allWords {
      combinedWord := fmt.Sprintf("%s%d", word, i)
      ipBlock := getIPBlock(combinedWord, isIPv6)
      if ipBlock == block {
        return combinedWord
      }
    }
  }
}

func isValidBlock(block string, isIPv6 bool) bool {
  if isIPv6 {
    return isValidHex(block)
  }
  return isValidNumber(block)
}

func isValidHex(block string) bool {
  if len(block) > 4 {
    return false
  }
  for _, char := range block {
    if !unicode.IsDigit(char) && (char < 'a' || char > 'f') {
      return false
    }
  }
  return true
}

func isValidNumber(block string) bool {
  if _, err := strconv.Atoi(block); err != nil {
    return false
  }
  return true
}

func ipType(isIPv6 bool) string {
  if isIPv6 {
    return "IPv6"
  }
  return "IPv4"
}

func blockType(isIPv6 bool) string {
  if isIPv6 {
    return "hex"
  }
  return "number"
}
