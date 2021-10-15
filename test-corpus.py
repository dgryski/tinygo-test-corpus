#!/usr/bin/python

import os
import sys

repos = [
    "dgryski/go-arc",
    "dgryski/go-bloomindex",
    "dgryski/go-camellia",
    "dgryski/go-change",
    "dgryski/go-chaskey",
    "dgryski/go-clefia",
    "dgryski/go-clockpro",
    "dgryski/go-cuckoof",
    "dgryski/go-discreterand",
    "dgryski/go-expirecache",
    "dgryski/go-factor",
    "dgryski/go-farm",
    "dgryski/go-fuzzstr",
    "dgryski/go-hollow",
    "dgryski/go-idea",
    "dgryski/go-interp",
    "dgryski/go-intpat",
    "dgryski/go-jump",
    "dgryski/go-kcipher2",
    "dgryski/go-ketama",
    "dgryski/go-krcrypt",
    "dgryski/go-linebreak",
    "dgryski/go-linlog",
    "dgryski/go-marvin32",
    "dgryski/go-md5crypt",
    "dgryski/go-metro",
    "dgryski/go-misty1",
    "dgryski/go-mph",
    "dgryski/go-neeva",
    "dgryski/go-nibz",
    "dgryski/go-pcgr",
    "dgryski/go-present",
    "dgryski/go-quicklz",
    "dgryski/go-radixsort",
    "dgryski/go-rbo",
    "dgryski/go-rc5",
    "dgryski/go-rc6",
    "dgryski/go-s4lru",
    "dgryski/go-sequitur",
    "dgryski/go-sip13",
    "dgryski/go-skinny",
    "dgryski/go-skip32",
    "dgryski/go-skipjack",
    "dgryski/go-sparx",
    "dgryski/go-spooky",
    "dgryski/go-spritz",
    "dgryski/go-timewindow",
    "dgryski/go-tinymap",
    "dgryski/go-twine",
    "dgryski/go-xoroshiro",
    "dgryski/go-xoshiro",

    # "dgryski/go-stablepart" -- requires reflect.DeepEqual() and testing/quick
    # "dgrysk/go-mavleg" -- requires reflect.DeepEqual
    # "dgryski/go-cobs", -- requires testing/quick
    # "dgryski/go-gramgen" -- requires building and running code and comparing output
    # "dgryski/go-kll", -- requires encoding/gob
    # "dgryski/go-mpchash", -- compat tests require siphash
    # "dgryski/go-nibblesort" -- requires testing/quick
    # "dgryski/go-postings" -- requires reflect.DeepEqual()
    # "dgryski/go-qselect" -- requires testing/quick
    # "dgryski/go-speck" -- requires reflect.DeepEqual()
    # "dgryski/go-trigram" -- requires reflect.DeepEqual()
    # "dgryski/go-simstore", -- requires testing/quick but can be moved to tinyfuzz with PR
    # "dgryski/go-ddmin" -- requires testing/quick
    # "dgryski/go-topk" -- requires encoding/gob
    # "dgryski/tsip/go" -- requires supporting cd'ing inside a repo
]


base_dir = os.getcwd()
corpus_dir = os.path.join(base_dir, "corpus")

# clone user/repo into user/repo directory, or update if it exists
def clone_or_update_repo(repo):
    if os.path.isdir(repo):
        os.chdir(repo)
        print("%s exists, updating" % repo)
        os.system("git fetch && git pull")
    else:
        print("%s does not exist, cloning" % repo)
        d = os.path.dirname(repo)
        if not os.path.isdir(d):
            print("creating directory %s" % d)
            os.makedirs(d)
        os.chdir(d)
        os.system("git clone https://github.com/%s" % repo)

def main():
    if not os.path.isdir(corpus_dir):
        os.makedirs(corpus_dir)


    for repo in repos:
        os.chdir(corpus_dir)
        clone_or_update_repo(repo)
        os.chdir(os.path.join(corpus_dir, repo))

        if not os.path.isfile("go.mod"):
            print ("creating running `go mod init")
            os.system("go mod init")
            os.system("go get -t .")

        for cmd in ("go test -v", "tinygo test -v -short -tags='purego noasm'"):
            print ("running `%s`" % cmd)
            r = os.system(cmd)
            if r:
                print ("`%s` failed" % cmd)
                sys.exit(1)


if __name__ == "__main__":
    main()
