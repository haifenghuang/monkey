fn generator() {
  ch = chan()
  n = 0
  spawn fn() {
    for {
      ch.send(n)
      n = n + 1
    }
    ch.close()
  }()
  return ch
}

for i in generator() {
  if i== 15 break
  fmt.println(i)
}

fmt.println("===========================")
fn generateUpperCaseLetters(inputString) {
    // Create a channel where to send output
    outputChannel = chan() 
    // Launch an (anonymous) function in another thread, that does
    // the actual processing.
    spawn fn() {
        // Loop over the letters in inputString
        for letter in inputString {
            // Send an uppercased letter to the output channel
            outputChannel.send(letter.upper())
        }
        // Close the output channel, so anything that loops over it
        // will know that it is finished.
        outputChannel.close()
    }()
    return outputChannel
}

// Loop over the letters communicated over the channel returned
// from generateUpperCaseLetters() and print them
for letter in generateUpperCaseLetters("hej") {
    fmt.println(letter)
}
