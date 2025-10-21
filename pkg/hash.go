package pkg

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"sort"
	"time"
)

func HashObject(content []byte) (string, bytes.Buffer, error) {
	// 1. Create a buffer to build the Git blob object.
	var buffer bytes.Buffer
	// 2. Write the blob header, including the type ("blob"), a space and the length of the content.
	// The `len()` function in Go returns the number of bytes in the slice.
	buffer.WriteString(fmt.Sprintf("blob %d", len(content)))
	// 3. Write the null byte ('\0'), which separates the header from the content.
	buffer.WriteByte(0)
	// 4. Write the actual content (the `[]byte`).
	buffer.Write(content)
	// 5. Compute the SHA-1 hash of the entire byte sequence.
	hash := sha1.Sum(buffer.Bytes())
	// 6. Format the resulting hash as a hexadecimal string.
	blobHash := fmt.Sprintf("%x", hash)
	return blobHash, buffer, nil
}

func HashCommit(parentHash, author, message string, files map[string]string) (string, []byte, error) {
	// 1. Use a buffer to efficiently build the commit content.
	var contentBuffer bytes.Buffer

	// 2. Write the commit metadata.
	// Fprintf is ideal for writing formatted text to an io.Writer like a buffer.
	fmt.Fprintf(&contentBuffer, "parent %s\n", parentHash)
	fmt.Fprintf(&contentBuffer, "author %s\n", author)
	// We use the ISO 8601 format (RFC3339 in Go) and UTC for consistency.
	fmt.Fprintf(&contentBuffer, "date %s\n", time.Now().UTC().Format(time.RFC3339))

	// 3. Write the commit message, separated by a blank line.
	fmt.Fprintf(&contentBuffer, "\n%s\n", message)

	// 4. Write the file list.
	if len(files) > 0 {
		// To ensure a deterministic commit hash, we must sort the files by their path.
		// Maps in Go do not guarantee iteration order.
		var paths []string
		for path := range files {
			paths = append(paths, path)
		}
		sort.Strings(paths) // Sort alphabetically.

		fmt.Fprintf(&contentBuffer, "\nfiles:\n")
		for _, path := range paths {
			hash := files[path]
			fmt.Fprintf(&contentBuffer, "%s %s\n", path, hash)
		}
	}

	// The commit content is ready.
	commitContent := contentBuffer.Bytes()

	// --- Now, we calculate the hash of the "commit object" Git-style ---
	// This is analogous to your HashObject function, but with the "commit" type.

	// We create a new buffer for the complete object (header + content).
	var objectBuffer bytes.Buffer
	// We write the header: "commit" type, a space, the content length, and a null byte.
	fmt.Fprintf(&objectBuffer, "commit %d\000", len(commitContent))
	// We write the content we just built.
	objectBuffer.Write(commitContent)

	// We calculate the SHA-1 hash of the complete object.
	hashBytes := sha1.Sum(objectBuffer.Bytes())
	commitHash := fmt.Sprintf("%x", hashBytes)

	// We return the commit hash and its content (without the "commit ..." header).
	return commitHash, commitContent, nil
}
