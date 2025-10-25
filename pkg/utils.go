package pkg

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetHeadRef() (map[string]string, error) {
	headRef := make(map[string]string)
	ref, err := os.Open(HeadPath)
	if err != nil {
		return nil, err
	}
	defer ref.Close()

	// 3. Create a scanner to read the file line by line
	scannerRef := bufio.NewScanner(ref)

	// 4. Iterate over each line of the file
	for scannerRef.Scan() {
		line := scannerRef.Text() // Get the line as a string

		// 5. Split the line into a slice of words
		words := strings.Fields(line)

		// 6. Check that there are at least two words
		if len(words) < 2 {
			log.Printf("Skipping line with incorrect format: %s", line)
			continue // Go to the next line if the format is incorrect
		}

		key := words[0]
		value := words[1]
		headRef[key] = value
	}

	// 8. Check for errors during scanning
	if err := scannerRef.Err(); err != nil {
		return nil, fmt.Errorf("error scanning HEAD file: %w", err)
	}

	return headRef, nil
}

// readIndex reads the index file into a map.
func ReadIndex() (map[string]string, error) {
	indexEntries := make(map[string]string)
	indexFile, err := os.Open(IndexPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return an empty map. It will be created on write.
			return indexEntries, nil
		}
		return nil, fmt.Errorf("error opening index for reading: %w", err)
	}
	defer indexFile.Close()

	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			indexEntries[parts[1]] = parts[0] // map[filepath] = hash
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning index file: %w", err)
	}
	return indexEntries, nil
}

// writeIndex writes the map of entries to the index file.
func WriteIndex(indexEntries map[string]string) error {
	var lines []string
	// For deterministic output, sort the file paths before writing.
	var paths []string
	for path := range indexEntries {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		lines = append(lines, fmt.Sprintf("%s %s", indexEntries[path], path))
	}

	output := strings.Join(lines, "\n")
	if len(lines) > 0 {
		output += "\n" // Add a final newline
	}

	if err := os.WriteFile(IndexPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("error writing to index file %s: %w", IndexPath, err)
	}
	return nil
}

func GetBranchHash() (string, error) {
	headFile, err := os.Open(HeadPath)
	if err != nil {
		return "", err
	}
	defer headFile.Close()

	var headRef string
	headScanner := bufio.NewScanner(headFile)
	for headScanner.Scan() {
		line := headScanner.Text()

		words := strings.Fields(line)
		headRef = words[1]
	}

	branchRefPath := fmt.Sprintf("%s/%s", RepoPath, headRef)
	branchHashFile, err := os.Open(branchRefPath)
	if err != nil {
		return "", err
	}
	defer branchHashFile.Close()

	var currentHash string
	brandHashScanner := bufio.NewScanner(branchHashFile)
	for brandHashScanner.Scan() {
		currentHash = brandHashScanner.Text()
	}

	return currentHash, nil
}

// BuildWorkdirMap recorre repoRoot y devuelve un map rutaRel -> sha1hex
func BuildWorkdirMap() (map[string]string, error) {
	repoRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("No se pudo obtener el directorio actual: %v", err)
	}
	workdirMap := make(map[string]string)

	// 1. Cargar las reglas de .gogitignore
	ignorePatterns, err := parseGitignore(repoRoot)
	if err != nil {
		return nil, err
	}

	// 2. Iniciar el recorrido recursivo
	walkErr := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Obtener la ruta relativa a la raíz del repositorio
		relativePath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		// Normalizar a forward slashes para comparación consistente
		relativePath = filepath.ToSlash(relativePath)

		// Saltar la raíz (".")
		if relativePath == "." {
			return nil
		}

		// Ignorar el propio archivo .gogitignore
		if !d.IsDir() && filepath.Base(relativePath) == ".gogitignore" {
			return nil
		}

		// Filtro Estricto: ignorar siempre el directorio .gogit
		if d.IsDir() && (relativePath == ".gogit" || strings.HasPrefix(relativePath, ".gogit/")) {
			return filepath.SkipDir
		}

		// Evaluar las reglas de ignore (se aplican en orden; las negaciones '!' deshacen ignores previos)
		isIgnored := false
		name := d.Name() // nombre del archivo/dir
		isDir := d.IsDir()

		for _, rawPattern := range ignorePatterns {
			if rawPattern == "" {
				continue
			}
			pattern := filepath.ToSlash(strings.TrimSpace(rawPattern))

			negated := false
			if strings.HasPrefix(pattern, "!") {
				negated = true
				pattern = strings.TrimPrefix(pattern, "!")
				pattern = strings.TrimSpace(pattern)
				if pattern == "" {
					// patrón "!" inválido -> ignorar
					continue
				}
			}

			// Si el pattern termina en "/" significa que apunta a directorios
			patternDirOnly := strings.HasSuffix(pattern, "/")
			if patternDirOnly {
				pattern = strings.TrimSuffix(pattern, "/")
			}

			matched := false

			// Si el pattern contiene una '/' hacemos la comparación contra la ruta relativa completa
			if strings.Contains(pattern, "/") {
				// si pattern empieza con "/" lo tratamos como relativo a la raíz: eliminamos prefijo si existe
				if strings.HasPrefix(pattern, "/") {
					pattern = strings.TrimPrefix(pattern, "/")
				}
				// Match usando filepath.Match contra relativePath
				if ok, matchErr := filepath.Match(pattern, relativePath); matchErr == nil && ok {
					matched = true
				} else if matchErr != nil {
					// patrón inválido — lo ignoramos
					continue
				}
			} else {
				// no contiene '/', comparar contra el nombre del archivo/dir
				if ok, matchErr := filepath.Match(pattern, name); matchErr == nil && ok {
					matched = true
				} else if matchErr != nil {
					continue
				}
			}

			// Si el patrón es exclusivo para directorios, y esto no es un dir -> no match
			if matched && patternDirOnly && !isDir {
				matched = false
			}

			if matched {
				if negated {
					// Una negación deshace el estado de ignore
					isIgnored = false
				} else {
					isIgnored = true
				}
				// No rompemos; git procesa todas las líneas (última coincidencia relevante)
			}
		}

		// Si está ignorado -> si es dir, evitar entrar; si es archivo, omitirlo.
		if isIgnored {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Si es un directorio válido no hacemos nada (solo archivos se hashean)
		if d.IsDir() {
			return nil
		}

		// 3. Procesar y hashear cada archivo válido
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("no se pudo leer el archivo %s: %w", path, err)
		}

		// Crear encabezado "blob <tamaño>\0"
		header := fmt.Sprintf("blob %d\x00", len(content))

		// Concatenar y hashear
		hasher := sha1.New()
		_, _ = hasher.Write([]byte(header))
		_, _ = hasher.Write(content)
		hashBytes := hasher.Sum(nil)

		// Convertir a hexadecimal y almacenar
		hashHex := hex.EncodeToString(hashBytes)
		// Guardar con la ruta relativa (sin "./")
		workdirMap[relativePath] = hashHex

		return nil
	})

	if walkErr != nil {
		return nil, fmt.Errorf("error durante el recorrido del directorio: %w", walkErr)
	}

	return workdirMap, nil
}

// parseGitignore lee .gogitignore y devuelve las líneas en orden (incluye negaciones)
func parseGitignore(repoRoot string) ([]string, error) {
	ignoreFilePath := filepath.Join(repoRoot, ".gogitignore")

	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		// Si no hay .gogitignore, no es un error, simplemente no hay nada que ignorar.
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error al leer .gogitignore: %w", err)
	}

	var patterns []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Ignorar líneas vacías y comentarios
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		patterns = append(patterns, trimmedLine)
	}
	return patterns, nil
}
