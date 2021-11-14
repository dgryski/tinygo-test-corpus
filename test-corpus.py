#!/usr/bin/python

import os
import sys

repos = [
    {
        'repo': 'dgryski/go-arc'
    },
    {
        'repo': 'dgryski/go-bloomindex',
        'tags': 'purego noasm',
    },
    {
        'repo': 'dgryski/go-camellia'
    },
    {
        'repo': 'dgryski/go-change'
    },
    {
        'repo': 'dgryski/go-chaskey',
        'tags': 'purego noasm',
    },
    {
        'repo': 'dgryski/go-clefia'
    },
    {
        'repo': 'dgryski/go-clockpro'
    },
    {
        'repo': 'dgryski/go-cobs'
    },
    {
        'repo': 'dgryski/go-cuckoof',
        'tags': 'pureno noasm',
    },
    {
        'repo': 'dgryski/go-discreterand'
    },
    {
        'repo': 'dgryski/go-expirecache'
    },
    {
        'repo': 'dgryski/go-factor'
    },
    {
        'repo': 'dgryski/go-farm',
        'tags': 'purego noasm',
    },
    {
        'repo': 'dgryski/go-fuzzstr'
    },
    {
        'repo': 'dgryski/go-hollow'
    },
    {
        'repo': 'dgryski/go-idea'
    },
    {
        'repo': 'dgryski/go-interp'
    },
    {
        'repo': 'dgryski/go-intpat'
    },
    {
        'repo': 'dgryski/go-jump'
    },
    {
        'repo': 'dgryski/go-kcipher2'
    },
    {
        'repo': 'dgryski/go-ketama'
    },
    {
        'repo': 'dgryski/go-krcrypt'
    },
    {
        'repo': 'dgryski/go-linebreak'
    },
    {
        'repo': 'dgryski/go-linlog'
    },
    {
        'repo': 'dgryski/go-maglev',
        'tags': 'appengine',  # for dchest/siphash
    },
    {
        'repo': 'dgryski/go-marvin32',
        'tags': 'purego'
    },
    {
        'repo': 'dgryski/go-md5crypt'
    },
    {
        'repo': 'dgryski/go-metro',
        'tags': 'purego noasm',
    },
    {
        'repo': 'dgryski/go-misty1'
    },
    {
        'repo': 'dgryski/go-mph',
        'tags': 'purego noasm',
    },
    {
        'repo': 'dgryski/go-neeva'
    },
    {
        'repo': 'dgryski/go-nibz'
    },
    {
        'repo': 'dgryski/go-nibblesort'
    },
    {
        'repo': 'dgryski/go-pcgr'
    },
    {
        'repo': 'dgryski/go-present'
    },
    {
        'repo': 'dgryski/go-quicklz'
    },
    {
        'repo': 'dgryski/go-radixsort'
    },
    {
        'repo': 'dgryski/go-rbo'
    },
    {
        'repo': 'dgryski/go-rc5'
    },
    {
        'repo': 'dgryski/go-rc6'
    },
    {
        'repo': 'dgryski/go-s4lru'
    },
    {
        'repo': 'dgryski/go-sequitur'
    },
    {
        'repo': 'dgryski/go-sip13',
        'tags': 'purego noasm',
    },
    {
        'repo': 'dgryski/go-skinny'
    },
    {
        'repo': 'dgryski/go-skip32'
    },
    {
        'repo': 'dgryski/go-skipjack'
    },
    {
        'repo': 'dgryski/go-sparx'
    },
    {
        'repo': 'dgryski/go-spooky'
    },
    {
        'repo': 'dgryski/go-spritz'
    },
    {
        'repo': 'dgryski/go-timewindow'
    },
    {
        'repo': 'dgryski/go-tinymap'
    },
    {
        'repo': "dgryski/go-trigram",
    },
    {
        'repo': 'dgryski/go-twine'
    },
    {
        'repo': 'dgryski/go-xoroshiro'
    },
    {
        'repo': 'dgryski/go-xoshiro'
    },
    {
        'repo': 'dgryski/go-zlatlong'
    },
    {
        'repo':
        'golang/crypto',
        'tags':
        'purego noasm',
        'subdirs': [
            'argon2',
            'bcrypt',
            'blake2b',
            'blake2s',
            'blowfish',
            'bn256',
            'cast5',
            'chacha20poly1305',
            'curve25519',
            'ed25519',
            'hkdf',
            'internal/subtle',
            'md4',
            'nacl/auth',
            'nacl/box',
            'nacl/secretbox',
            'nacl/sign',
            'openpgp/armor',
            'openpgp/elgamal',
            'openpgp/s2k',
            'pbkdf2',
            'pkcs12/internal/rc2',
            'ripemd160',
            'salsa20',
            'scrypt',
            'ssh/internal/bcrypt_pbkdf',
            'tea',
            'twofish',
            'xtea',
            # chacha20 -- panic: chacha20: SetCounter attempted to rollback counter
            # cryptobyte -- panic: unimplemented: reflect.OverflowInt()
            # salsa20/salsa -- panic: runtime error: index out of range
        ]
    },
    {
        'repo': 'google/shlex'
    },
    {
        'repo': 'google/boundedwait'
    },
    {
        'repo': 'dgryski/go-maglev',
        'tags': 'appengine',  # for dchest/siphash
    },
    {
        'repo': 'google/btree',
    },
    {
        'repo': 'google/der-ascii',
        'subdirs': ['cmd/ascii2der', 'cmd/der2ascii', 'internal'],
    },
    {
        'repo': 'google/hilbert'
    },
    {
        'repo': 'google/go-intervals',
        'subdirs': ['intervalset', 'timespanset'],
    },
    {
        'repo': 'google/okay'
    },
    {
        'repo':
        'golang/text',
        'subdirs': [
            'encoding',
            'encoding/charmap',
            'encoding/htmlindex',
            'encoding/ianaindex',
            'encoding/japanese',
            'encoding/korean',
            'encoding/simplifiedchinese',
            'encoding/traditionalchinese',
            'encoding/unicode',
            'encoding/utf32',
            'internal/format',
            "internal/ucd",
            "internal/tag",
            'search',
            'unicode/rangetable',
            # internal/stringset -- fails due to sort.Search()?
        ]
    },
    {
        'repo':
        'golang/image',
        'subdirs': [
            'ccitt',
            'colornames',
            'draw',
            'font',
            'font/basicfont',
            'font/opentype',
            'font/plan9font',
            'math/fixed',
            'riff',
            # 'tiff', -- fails "panic: runtime error: nil pointer dereference"
            'webp',
        ],
        'tags':
        'noasm',
    },
    {
        'repo':
        'golang/geo',
        'subdirs': [
            'r1',
            'r2',
            'r3',
            's1',
            #  's2', -- fails, possibly due to sort.Search() bug
        ],
    },
    {
        'repo': 'golang/groupcache',
        'subdirs': [
            'consistenthash',
            'lru',
        ],
    },
    # "dgryski/go-stablepart" -- requires reflect.DeepEqual() and testing/quick
    # "dgryski/go-cobs", -- requires testing/quick
    # "dgryski/go-gramgen" -- requires building and running code and comparing output
    # "dgryski/go-kll", -- requires encoding/gob
    # "dgryski/go-mpchash", -- compat tests require siphash
    # "dgryski/go-qselect" -- requires testing/quick
    # "dgryski/go-simstore", -- requires testing/quick but can be moved to tinyfuzz with PR
    # "dgryski/go-ddmin" -- requires testing/quick
    # "dgryski/go-topk" -- requires encoding/gob
    # "golang/snappy" -- needs patching out os.* bits; target=wasi hangs?
    # "cloudflare/ahocorasick" -- interp timeout building regexps in test
    # "google/open-location-code/go" -- alloc link error
    # "dgryski/go-postings" -- nil map in AddDocument causes segfault
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
    os.system("tinygo clean")
    if not os.path.isdir(corpus_dir):
        os.makedirs(corpus_dir)

    for repo in repos:
        os.chdir(corpus_dir)
        clone_or_update_repo(repo['repo'])
        repo_base = os.path.join(corpus_dir, repo['repo'])
        os.chdir(repo_base)

        if not os.path.isfile("go.mod"):
            print("creating go.mod: running `go mod init")
            os.system("go mod init github.com/%s" % repo['repo'])
            os.system("go get -t .")

        tags = ""
        if 'tags' in repo:
            tags = "-tags='%s'" % repo['tags']

        for cmd in ["go test -v", "tinygo test -v %s" % tags]:
            dirs = ["."]
            if 'subdirs' in repo:
                dirs = repo['subdirs']

            for subdir in dirs:
                if subdir != ".":
                    os.chdir(subdir)
                print("running `%s`" % cmd)
                r = os.system(cmd)
                if r:
                    print("`%s` failed" % cmd)
                    sys.exit(1)
                if subdir != ".":
                    os.chdir(repo_base)


if __name__ == "__main__":
    main()
