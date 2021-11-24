package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	// underscore prefix so go tool excludes corpus directory.
	corpusFolderName = "_corpus"
	dirMode          = 0777
	host             = "github.com"
	hostURL          = "https://" + host
)

func main() {
	var countSubdir, countRepo int
	defer func() {
		log.Printf("Finished!\n%d/%d repos tested\n%d passed subdir tests\n", countRepo, len(repos), countSubdir)
	}()

	// Workspace setup and cleanup.
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal("getting current dir:", err)
	}
	corpusDir := filepath.Join(baseDir, corpusFolderName)
	mustrun("tinygo", "clean")
	if err != nil {
		log.Fatal("calling `tinygo clean`:", err)
	}
	os.Mkdir(corpusDir, dirMode) // force directory creation if not exist.
	_, err = os.ReadDir(corpusDir)
	if err != nil {
		log.Fatal("reading corpus directory: ", err)
	}

	// Commence testing logic.
	for _, repo := range repos {
		os.Chdir(corpusDir)
		cloneOrUpdateRepo(repo.Repo)
		repoBase := filepath.Join(corpusDir, repo.Repo)
		os.Chdir(repoBase)

		if _, err := os.Stat("go.mod"); err != nil {
			log.Printf("creating %s/go.mod: running `go mod init`\n", repoBase)
			mustrun("go", "mod", "init", fmt.Sprintf("%s/%s", host, repo.Repo))
			mustrun("go", "get", "-t", ".")
		}
		tags := ""
		if repo.Tags != "" {
			tags = fmt.Sprintf("%s", repo.Tags)
		}
		dirs := []string{"."}
		if len(repo.Subdirs) > 0 {
			dirs = repo.Subdirs
		}

		for _, subdir := range dirs {
			if subdir != "." {
				os.Chdir(subdir)
			}
			tinyout := make(chan string)
			// Run TinyGo and Go in parallel.
			go func() {
				tinyout <- mustrun("tinygo", "test", "-v", "-tags", tags)
			}()
			out1 := mustrun("go", "test", "-v")
			countSubdir++
			log.Printf("package %s:\n%s\n%s\n", filepath.Join(repo.Repo, subdir), out1, <-tinyout)
			if subdir != "." {
				os.Chdir(repoBase)
			}
		}
		countRepo++
		log.Printf("finished module %d/%d %s", countRepo, len(repos), repo.Repo)
	}
}

func cloneOrUpdateRepo(repo string) {
	if _, err := os.Stat(repo); err != nil {
		// Repo does not exist.
		log.Printf("repo not found. cloning %s", repo)
		d := filepath.Dir(repo)
		if _, err := os.Stat(repo); err != nil {
			log.Printf("creating directory %s", d)
			os.Mkdir(d, dirMode)
		}
		os.Chdir(d)
		mustrun("git", "clone", fmt.Sprintf("%s/%s", hostURL, repo))
		return
	}

	os.Chdir(repo)
	log.Printf("repo exists, updating %s", repo)
	mustrun("git", "fetch")
	mustrun("git", "pull")
}

func mustrun(name string, arg ...string) (stdout string) {
	cmd := exec.Command(name, arg...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		cwd, _ := os.Getwd()
		log.Fatalf("%s\ncmd %s with err: %q at dir %q", string(b), cmd.String(), err, cwd)
	}
	return string(b)
}

type T struct {
	Repo    string
	Tags    string
	Subdirs []string
}

var repos = []T{
	{
		Repo: "dgryski/go-bloomindex",
		Tags: "purego noasm",
	},
	{
		Repo: "dgryski/go-arc",
	},
	{
		Repo: "dgryski/go-camellia",
	},
	{
		Repo: "dgryski/go-change",
	},
	{
		Repo: "dgryski/go-chaskey",
		Tags: "purego noasm",
	},
	{
		Repo: "dgryski/go-clefia",
	},
	{
		Repo: "dgryski/go-clockpro",
	},
	{
		Repo: "dgryski/go-cobs",
	},
	{
		Repo: "dgryski/go-cuckoof",
		Tags: "pureno noasm",
	},
	{
		Repo: "dgryski/go-discreterand",
	},
	{
		Repo: "dgryski/go-expirecache",
	},
	{
		Repo: "dgryski/go-factor",
	},
	{
		Repo: "dgryski/go-farm",
		Tags: "purego noasm",
	},
	{
		Repo: "dgryski/go-fuzzstr",
	},
	{
		Repo: "dgryski/go-hollow",
	},
	{
		Repo: "dgryski/go-idea",
	},
	{
		Repo: "dgryski/go-interp",
	},
	{
		Repo: "dgryski/go-intpat",
	},
	{
		Repo: "dgryski/go-jump",
	},
	{
		Repo: "dgryski/go-kcipher2",
	},
	{
		Repo: "dgryski/go-ketama",
	},
	{
		Repo: "dgryski/go-krcrypt",
	},
	{
		Repo: "dgryski/go-linebreak",
	},
	{
		Repo: "dgryski/go-linlog",
	},
	{
		Repo: "dgryski/go-maglev",
		Tags: "appengine", // for dchest/siphash
	},
	{
		Repo: "dgryski/go-marvin32",
		Tags: "purego",
	},
	{
		Repo: "dgryski/go-md5crypt",
	},
	{
		Repo: "dgryski/go-metro",
		Tags: "purego noasm",
	},
	{
		Repo: "dgryski/go-misty1",
	},
	{
		Repo: "dgryski/go-mph",
		Tags: "purego noasm",
	},
	{
		Repo: "dgryski/go-mpchash",
		Tags: "appengine", // for dchest/siphash
	},
	{
		Repo: "dgryski/go-neeva",
	},
	{
		Repo: "dgryski/go-nibz",
	},
	{
		Repo: "dgryski/go-nibblesort",
	},
	{
		Repo: "dgryski/go-pcgr",
	},
	{
		Repo: "dgryski/go-present",
	},
	{
		Repo: "dgryski/go-quicklz",
	},
	{
		Repo: "dgryski/go-radixsort",
	},
	{
		Repo: "dgryski/go-rbo",
	},
	{
		Repo: "dgryski/go-rc5",
	},
	{
		Repo: "dgryski/go-rc6",
	},
	{
		Repo: "dgryski/go-s4lru",
	},
	{
		Repo: "dgryski/go-sequitur",
	},
	{
		Repo: "dgryski/go-sip13",
		Tags: "purego noasm",
	},
	{
		Repo: "dgryski/go-skinny",
	},
	{
		Repo: "dgryski/go-skip32",
	},
	{
		Repo: "dgryski/go-skipjack",
	},
	{
		Repo: "dgryski/go-sparx",
	},
	{
		Repo: "dgryski/go-spooky",
	},
	{
		Repo: "dgryski/go-spritz",
	},
	{
		Repo: "dgryski/go-timewindow",
	},
	{
		Repo: "dgryski/go-tinymap",
	},
	{
		Repo: "dgryski/go-trigram",
	},
	{
		Repo: "dgryski/go-twine",
	},
	{
		Repo: "dgryski/go-xoroshiro",
	},
	{
		Repo: "dgryski/go-xoshiro",
	},
	{
		Repo: "dgryski/go-zlatlong",
	},
	{
		Repo: "golang/crypto",
		Tags: "purego noasm",
		Subdirs: []string{
			"argon2",
			"bcrypt",
			"blake2b",
			"blake2s",
			"blowfish",
			"bn256",
			"cast5",
			"chacha20poly1305",
			"curve25519",
			"ed25519",
			"hkdf",
			"internal/subtle",
			"md4",
			"nacl/auth",
			"nacl/box",
			"nacl/secretbox",
			"nacl/sign",
			"openpgp/armor",
			"openpgp/elgamal",
			"openpgp/s2k",
			"pbkdf2",
			"pkcs12/internal/rc2",
			"ripemd160",
			"salsa20",
			"scrypt",
			"ssh/internal/bcrypt_pbkdf",
			"tea",
			"twofish",
			"xtea",
			// chacha20 -- panic: chacha20: SetCounter attempted to rollback counter
			// cryptobyte -- panic: unimplemented: reflect.OverflowInt()
			// salsa20/salsa -- panic: runtime error: index out of range
		},
	},
	{
		Repo: "google/shlex",
	},
	{
		Repo: "google/boundedwait",
	},
	{
		Repo: "dgryski/go-maglev",
		Tags: "appengine", // for dchest/siphash
	},
	{
		Repo: "google/btree",
	},
	{
		Repo: "google/der-ascii",
		Subdirs: []string{
			"cmd/ascii2der",
			"cmd/der2ascii",
			"internal",
		},
	},
	{
		Repo: "google/hilbert",
	},
	{
		Repo: "google/go-intervals",
		Subdirs: []string{
			"intervalset",
			"timespanset",
		},
	},
	{
		Repo: "google/okay",
	},
	{
		Repo: "golang/text",
		Subdirs: []string{
			"encoding",
			"encoding/charmap",
			"encoding/htmlindex",
			"encoding/ianaindex",
			"encoding/japanese",
			"encoding/korean",
			"encoding/simplifiedchinese",
			"encoding/traditionalchinese",
			"encoding/unicode",
			"encoding/utf32",
			"internal/format",
			"internal/ucd",
			"internal/tag",
			"search",
			"unicode/rangetable",
			// internal/stringset -- fails due to sort.Search()?
		},
	},
	{
		Repo: "golang/image",
		Tags: "noasm",
		Subdirs: []string{
			"ccitt",
			"colornames",
			"draw",
			"font",
			"font/basicfont",
			"font/opentype",
			"font/plan9font",
			"math/fixed",
			"riff",
			//  "tiff", -- fails "panic: runtime error: nil pointer dereference"
			"webp",
		},
	},
	{
		Repo: "golang/geo",
		Subdirs: []string{
			"r1",
			"r2",
			"r3",
			"s1",
			// "s2", -- fails, possibly due to sort.Search() bug
		},
	},
	{
		Repo: "golang/groupcache",
		Subdirs: []string{
			"consistenthash",
			"lru",
		},
	},
	{
		Repo: "armon/go-radix",
	},
	{
		Repo: "armon/circbuf",
	},
	{
		Repo: "VividCortex/gohistogram",
	},
	{
		Repo: "cespare/xxhash",
		Tags: "appengine",
	},
	{
		Repo: "gonum/gonum",
		Tags: "noasm",
	},
	{
		Repo: "gonum/gonum",
		Tags: "noasm",
		Subdirs: []string{
			"blas/blas32",
			// "blas/blas64", -- TestDasum panic: blas: n < 0
			"blas/cblas64",
			"blas/cblas128",
			// "blas/gonum", -- panic: blas: n < 0
			// "cmplxs", -- TestAdd panic: cmplxs: slice lengths do not match
			"cmplxs/cscalar",
			// "diff/fd" -- panic: fd: slice length mismatch
			// "dsp/fourier" -- panic: unimplemented: reflect.DeepEqual()
			"dsp/window",
			// "floats", -- panic: floats: destination slice length does not match input
			"floats/scalar",
			// "graph" ld.lld-11:  -- error: undefined symbol: reflect.mapiterkey (among other reflect errors)
			// "graph/topo" -- Reflect: Same as above
			"integrate",
			"integrate/quad",
			"internal/cmplx64",
			// "internal/math32" -- /usr/local/go/src/testing/quick/quick.go:273:11: fType.NumOut undefined (type reflect.Type has no field or method NumOut)
			"internal/testrand",
			// "interp", -- panic: interp: input slices have different lengths
			"lapack/gonum",
			// "mat", -- panic: unimplemented: reflect.DeepEqual()
			"mathext",
			"mathext/prng",
			// "num/dual" -- TestFormat unexpected result for fmt.Sprintf("%//v", T{Real:1.1, Emag:2.1}): got:"T{Real:1.1, Emag:2.1}", want:"dual.Number{Real:1.1, Emag:2.1}"    unexpected result for fmt.Sprintf("%//v", T{Real:-1.1, Emag:-2.1}): got:"T{Real:-1.1, Emag:-2.1}", want:"dual.Number{Real:-1.1, Emag:-2.1}"
			// "num/dualcmplx" -- TestFormat (similar to above)
			// "num/dualquat" -- TestFormat (similar to above)
			// "num/hyperdual" -- TestFormat (similar to above)
			// "num/quat" -- TestFormat (similar to above)
			// "optimize" // ld.lld-11: error: undefined symbol: golang.org/x/tools/container/intsets.havePOPCNT error: failed to link ...
			"optimize/convex/lp",
			"optimize/functions",
			// "spatial/barneshut", -- panic: unimplemented: reflect.DeepEqual()
			// "spatial/kdtree", -- panic: unimplemented: reflect.DeepEqual()
			"spatial/r2",
			"spatial/r3",
			// "spatial/vptree", -- panic: unimplemented: reflect.DeepEqual()
			// "stat", -- panic: stat: slice length mismatch
			// "stat/card" -- /usr/local/go/src/encoding/gob/decode.go:562:21: MakeMapWithSize not declared by package reflect
			// "stat/combin" -- panic: unimplemented: reflect.DeepEqual()
			"stat/distmat",
			// "stat/distmv", -- ld.lld-11: error: undefined symbol: golang.org/x/tools/container/intsets.havePOPCNT error: failed to link ...
			// "stat/distuv" -- panic: distuv: cannot compute Mode for Beta != 0\
			"stat/mds",
			// "stat/samplemv", -- ld.lld-11: error: undefined symbol: golang.org/x/tools/container/intsets.havePOPCNT error: failed to link ...
			// "stat/sampleuv", -- panic: unimplemented: reflect.DeepEqual()
			"stat/spatial",
			// "unit" -- All Format tests fail. Similar to `num` subpackages.
		},

		// "dgryski/go-stablepart" -- requires reflect.DeepEqual() and testing/quick
		// "dgryski/go-cobs", -- requires testing/quick
		// "dgryski/go-gramgen" -- requires building and running code and comparing output
		// "dgryski/go-kll", -- requires encoding/gob
		// "dgryski/go-qselect" -- requires testing/quick
		// "dgryski/go-simstore", -- requires testing/quick but can be moved to tinyfuzz with PR
		// "dgryski/go-ddmin" -- requires testing/quick
		// "dgryski/go-topk" -- requires encoding/gob
		// "golang/snappy" -- needs patching out os.* bits; target=wasi hangs?
		// "cloudflare/ahocorasick" -- interp timeout building regexps in test
		// "google/open-location-code/go" -- alloc link error
		// "dgryski/go-postings" -- nil map in AddDocument causes segfault
	},
}
