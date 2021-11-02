#!/usr/bin/python

import os
import sys

repos = [
    {
        'repo': 'dgryski/go-arc'
    },
    {
        'repo': 'dgryski/go-bloomindex'
    },
    {
        'repo': 'dgryski/go-camellia'
    },
    {
        'repo': 'dgryski/go-change'
    },
    {
        'repo': 'dgryski/go-chaskey'
    },
    {
        'repo': 'dgryski/go-clefia'
    },
    {
        'repo': 'dgryski/go-clockpro'
    },
    {
        'repo': 'dgryski/go-cuckoof'
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
        'repo': 'dgryski/go-farm'
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
        'repo': 'dgryski/go-marvin32',
        'tags': 'purego'
    },
    {
        'repo': 'dgryski/go-md5crypt'
    },
    {
        'repo': 'dgryski/go-metro'
    },
    {
        'repo': 'dgryski/go-misty1'
    },
    {
        'repo': 'dgryski/go-mph'
    },
    {
        'repo': 'dgryski/go-neeva'
    },
    {
        'repo': 'dgryski/go-nibz'
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
        'repo': 'dgryski/go-sip13'
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
        'repo': 'dgryski/go-twine'
    },
    {
        'repo': 'dgryski/go-xoroshiro'
    },
    {
        'repo': 'dgryski/go-xoshiro'
    },
    {
        'repo':
        'golang/crypto',
        'tags':
        'purego noasm',
        'subdirs': [
            'argon2', 'bcrypt', 'blake2b', 'blake2s', 'blowfish', 'bn256',
            'cast5', 'chacha20poly1305', 'curve25519', 'ed25519', 'hkdf',
            'internal/subtle', 'md4', 'nacl/box', 'nacl/secretbox',
            'nacl/sign', 'openpgp/armor', 'openpgp/elgamal', 'openpgp/s2k',
            'pbkdf2', 'pkcs12/internal/rc2', 'ripemd160', 'salsa20', 'scrypt',
            'tea', 'twofish', 'xtea'
        ]
    },
    {
        'repo': 'jedisct1/go-minisign'
    },
    {
        'repo': 'jedisct1/xsecretbox',
        'tags': 'purego noasm'
    },
    {
        'repo': 'google/shlex'
    },

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
        clone_or_update_repo(repo['repo'])
        repo_base = os.path.join(corpus_dir, repo['repo'])
        os.chdir(repo_base)

        if not os.path.isfile("go.mod"):
            print("creating running `go mod init")
            os.system("go mod init github.com/%s" % repo)
            os.system("go get -t .")

        tags = ""
        if 'tags' in repo:
            tags = "-tags='%s'" % repo['tags']

        for cmd in ["go test -v", "tinygo test -v -short %s" % tags]:
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
