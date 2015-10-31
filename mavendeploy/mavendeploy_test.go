package mavendeploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkip(t *testing.T) {
	t.Parallel()
	mvn := &Maven{
		Repository: Repository{},
		Artifact:   Artifact{},
		GPG:        GPG{},
		Args:       Args{},
	}
	localTest(mvn, func() {
		err := mvn.Publish()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestURL(t *testing.T) {
	t.Parallel()
	mvn := &Maven{
		Repository: Repository{Username: "u", Password: "p"},
		Artifact:   Artifact{},
		GPG:        GPG{},
		Args:       Args{},
	}

	err := mvn.Publish()
	if err == nil || err != errRequiredValue {
		t.Fatal("url should be required")
	}
}

func TestPublish1(t *testing.T) {
	t.Parallel()
	mvn := &Maven{
		Repository: Repository{
			Username: "u",
			Password: "p",
		},
		Artifact: Artifact{
			GroupID: "com.test.publish1",
		},
		GPG: GPG{},
		Args: Args{
			Source: "multiple-matched/app*",
			Regexp: "(?P<artifact>app-[^/-]*)-(?P<classifier>[^-]*-[^-]*)-(?P<version>.*).(?P<extension>tar.gz|zip|readme)$",
		},
	}
	localTest(mvn, func() {
		err := mvn.Publish()
		if err != nil {
			t.Fatal(err)
		}
		assertLocalArtifacts(t, mvn,
			"com/test/publish1/app-client/maven-metadata.xml.md5",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.sha1",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.sha1",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.md5",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4.pom.md5",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.md5",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-windows-386.zip.sha1",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.sha1",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4.pom",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.sha1",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4.pom.sha1",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-windows-386.zip.md5",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.md5",
			"com/test/publish1/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.md5",
			"com/test/publish1/app-client/maven-metadata.xml.sha1",
			"com/test/publish1/app-client/maven-metadata.xml",
			"com/test/publish1/app-server/maven-metadata.xml.md5",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4-linux-amd64.readme.sha1",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4-linux-amd64.tar.gz.sha1",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4-linux-amd64.tar.gz.md5",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4-linux-amd64.readme.md5",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4.pom.md5",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4.pom",
			"com/test/publish1/app-server/0.1.4/app-server-0.1.4.pom.sha1",
			"com/test/publish1/app-server/maven-metadata.xml.sha1",
			"com/test/publish1/app-server/maven-metadata.xml",
			"com/test/publish1/app-gui/maven-metadata.xml.md5",
			"com/test/publish1/app-gui/0.1.4/app-gui-0.1.4-darwin-amd64.zip.md5",
			"com/test/publish1/app-gui/0.1.4/app-gui-0.1.4.pom",
			"com/test/publish1/app-gui/0.1.4/app-gui-0.1.4-darwin-amd64.zip.sha1",
			"com/test/publish1/app-gui/0.1.4/app-gui-0.1.4.pom.md5",
			"com/test/publish1/app-gui/0.1.4/app-gui-0.1.4.pom.sha1",
			"com/test/publish1/app-gui/maven-metadata.xml.sha1",
			"com/test/publish1/app-gui/maven-metadata.xml",
		)
	})
}

func TestPublish2(t *testing.T) {
	t.Parallel()

	mvn := &Maven{
		Repository: Repository{
			Username: "u",
			Password: "p",
		},
		Artifact: Artifact{
			GroupID:    "com.test.publish2",
			ArtifactID: "release",
			Extension:  "zip",
			Version:    "1.2.3",
		},
		GPG: GPG{},
		Args: Args{
			Source: "single/release.zip",
		},
	}
	localTest(mvn, func() {
		err := mvn.Publish()
		if err != nil {
			t.Fatal(err)
		}
		assertLocalArtifacts(t, mvn,
			"com/test/publish2/release/maven-metadata.xml.md5",
			"com/test/publish2/release/1.2.3/release-1.2.3.zip.md5",
			"com/test/publish2/release/1.2.3/release-1.2.3.zip.sha1",
			"com/test/publish2/release/1.2.3/release-1.2.3.pom.sha1",
			"com/test/publish2/release/1.2.3/release-1.2.3.pom",
			"com/test/publish2/release/1.2.3/release-1.2.3.pom.md5",
			"com/test/publish2/release/maven-metadata.xml.sha1",
			"com/test/publish2/release/maven-metadata.xml",
		)
	})
}

func TestPublish3(t *testing.T) {
	t.Parallel()

	mvn := &Maven{
		Repository: Repository{
			Username: "u",
			Password: "p",
		},
		Artifact: Artifact{
			GroupID:    "com.test.publish3",
			Extension:  "zip",
			Version:    "1.2.3",
			ArtifactID: "asd",
		},
		GPG: GPG{},
		Args: Args{
			Source: "single-matched/*.zip",
			Regexp: "(?P<artifact>[^/-]*)-(?P<classifier>[^-]*-[^-]*).zip$",
		},
	}
	localTest(mvn, func() {
		err := mvn.Publish()
		if err != nil {
			t.Fatal(err)
		}
		assertLocalArtifacts(t, mvn,
			"com/test/publish3/app/maven-metadata.xml.md5",
			"com/test/publish3/app/1.2.3/app-1.2.3.pom.md5",
			"com/test/publish3/app/1.2.3/app-1.2.3-windows-amd64.zip.sha1",
			"com/test/publish3/app/1.2.3/app-1.2.3.pom.sha1",
			"com/test/publish3/app/1.2.3/app-1.2.3.pom",
			"com/test/publish3/app/1.2.3/app-1.2.3-windows-amd64.zip.md5",
			"com/test/publish3/app/maven-metadata.xml.sha1",
			"com/test/publish3/app/maven-metadata.xml",
		)
	})

}

func TestGPGSign(t *testing.T) {
	t.Parallel()

	mvn := &Maven{
		Repository: Repository{
			Username: "user",
			Password: "pass",
		},
		Artifact: Artifact{
			GroupID:    "com.test.publishGpg",
			ArtifactID: "release",
			Extension:  "zip",
			Version:    "1.9.3",
		},
		GPG: GPG{
			PrivateKey: privateKey,
			Passphrase: `test`,
		},
		Args: Args{
			Source: "single/release.zip",
			Debug:  true,
		},
	}
	localTest(mvn, func() {
		err := mvn.Publish()
		if err != nil {
			t.Fatal(err)
		}
		assertLocalArtifacts(t, mvn,
			"com/test/publishGpg/release/maven-metadata.xml.md5",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.zip.asc.md5",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.pom",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.pom.asc",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.pom.asc.sha1",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.zip.md5",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.pom.asc.md5",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.pom.sha1",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.zip.sha1",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.pom.md5",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.zip.asc.sha1",
			"com/test/publishGpg/release/1.9.3/release-1.9.3.zip.asc",
			"com/test/publishGpg/release/maven-metadata.xml.sha1",
			"com/test/publishGpg/release/maven-metadata.xml",
		)
	})
}

func TestGPGSign2(t *testing.T) {
	t.Parallel()
	mvn := &Maven{
		Repository: Repository{
			Username: "u",
			Password: "p",
		},
		Artifact: Artifact{
			GroupID: "com.test.gpg2",
		},
		GPG: GPG{
			PrivateKey: privateKey,
			Passphrase: `test`,
		},
		Args: Args{
			Source: "multiple-matched/app-client*",
			Regexp: "(?P<artifact>app-[^/-]*)-(?P<classifier>[^-]*-[^-]*)-(?P<version>.*).(?P<extension>tar.gz|zip|readme)$",
		},
	}
	localTest(mvn, func() {
		err := mvn.Publish()
		if err != nil {
			t.Fatal(err)
		}

		assertLocalArtifacts(t, mvn,
			"com/test/gpg2/app-client/maven-metadata.xml.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4.pom.asc.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.asc",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.asc.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-386.zip.asc.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.asc.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4.pom.asc",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4.pom.asc.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-386.zip.asc",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.asc",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-386.zip.asc.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4.pom.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.asc.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-386.zip.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.asc.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-386.tar.gz.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4.pom",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.asc.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4.pom.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.asc.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-386.zip.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.asc",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.asc.sha1",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-darwin-amd64.zip.asc",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-windows-amd64.zip.asc.md5",
			"com/test/gpg2/app-client/0.1.4/app-client-0.1.4-linux-amd64.tar.gz.md5",
			"com/test/gpg2/app-client/maven-metadata.xml.sha1",
			"com/test/gpg2/app-client/maven-metadata.xml",
		)
	})
}

func localTest(mvn *Maven, fn func()) {
	tmpdir, err := ioutil.TempDir("", "drone-mvn-test")
	if err != nil {
		panic(err)
	}
	mvn.Repository.URL = fmt.Sprintf("file://%s", tmpdir)
	defer func() {
		os.RemoveAll(tmpdir)
	}()
	mvn.workspacePath = "test-data/"
	mvn.quiet = true
	fn()
}

func assertLocalArtifacts(t *testing.T, mvn *Maven, path ...string) {
	basepath := strings.TrimPrefix(mvn.Repository.URL, "file://")
	for _, v := range path {
		fullpath := filepath.Join(basepath, v)
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			t.Fatalf("no such file or directory: %s", fullpath)
			return
		}
	}

}

const privateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1

lQH+BFY0AN8BBADCJ7NAMFJXkgti6vpxCZSlZlO6IjqrEmHBnyLkIo6OX1uZmtBS
f1wlSVAevcNJJJaHkLQz8vAvE7lzxVvHEL8P2eg6zUGmJRElCbdcP6HtivYguat+
VdUe057Kp7sOyjhi7P2oUTUj7Ma3RGAvoi97uggBl88gwYLy+hz8MBPelQARAQAB
/gMDAg9KGIUVlIkuYNqxNsdnk++EHjebW/ONdwCuB7pPW0NKoBs7ekBqwKYor7KD
4JCgKY98e7FF8gJbDu0272x8WFgf9Svh9P1td9IPLiWomJJh+/KyhpkGgQXbC9XI
qQbyTiZXVsu+y/0SIKHbMUjh/AjaLbKgSjUu8sY8B53xnyzQ5wZkDwMtcRDIR5qi
niAjUP0nUt+WBA8mzJJKmR5qe2bjACw01sc32BYkGeopAbR8GtQVobowm4IXraU0
t2cZPfVU+kbRffljcoJw1IQGHY64QoNuxc666HWZBVi+Sw7l4xWvrE0gj1GDmXQj
yemwiRb00xBpih/G/Ha4l5ixWysuN9on6xU1KZI9Hikcz3BaNRoRFvfwcU9zXvUE
3ul9iqVy8Kbwa2fagjzdPmLSViru7KaqcQVehpL6OMKZM/GzvffWGrCSFvyevMIh
7191OmnmV7Nm5rmyNIhGRiUL0sp3KR/oVLbDB+FDfHtRtAh0ZXN0IGtleYi4BBMB
AgAiBQJWNADfAhsDBgsJCAcDAgYVCAIJCgsEFgIDAQIeAQIXgAAKCRCJ5RVIH4+h
LuTVA/9xMKoBLPuneU9ZpIbb5dBAnnnrDECKMxGF/9c+sIyfWF5vSumyIrB6VFMA
6iN0blbIBXacBncSTr5pW5eInqpB8Cs8FdiPyBiWhB6SGZXQarKS3cSZHk97bTN3
oBoH67kPlKnD4F+INqsj+em0iOmn3VwtaYepTHSdz24dcSFJDZ0B/gRWNADfAQQA
skGna66JAiAw7lQTYXnWqQ8Fw4tR5jRbXCSP3Sg0Yf/Y84cvHAQwUJDUlDdqqzqx
/Yr4NcyEJ8Kdux601aA9UhBDFIuoQQep6ETUnRzwqRWQmK/hT8L49wrmRqkjKxqR
OFgKDK0O1vHnAlh9kZc12XjjPDWB7l2EiXK5kgLGpesAEQEAAf4DAwIPShiFFZSJ
LmDCLRhxDFymfUypuHNkYEFj03+D4hpY7PAMpRSO++oP+psS1Y4DbdA+b96VR8xA
MK7p30HG2M829z9I8j9+HbhXqAXrnFqWQqf5XRmgcxaIWQyte7ZBa1nVFQN1fWiC
gYD+Uhlo6AauaKnxIqkZWog6QNat84QR3tSywfWmI91Avluhcqtp4oBjN/SR2m3R
XHaOCWNikG733CVv8ZxxwWcgZ4iEPDwrLEXs2W19ehygpJX50Z3n+85fKIsp2cGh
cLM6dlwZrlHzhRUy7NhOlmQaCNygW/kLzBO3uHEI5qElp+QhTxcgf3s72IaX4bgK
QAQ9BtLVLxiJop/mtFTgF3g9Fpxr3xe1LtgUTbnS0OIMiAst6Z/cbCKGSsl5Nl5I
WcWRPJEs6+Lx90nYHijrZt8/G27CwEN2UiqxE5dyccleCIUyzQ/KvwyjxS/BZm/+
rjs4nvUB0yxr3iqFlqKOO7uvjltkIYifBBgBAgAJBQJWNADfAhsMAAoJEInlFUgf
j6EuxdUEAMCnHTvReIvWkNKyzjzK5a0DZCmJLoFmJ8zmNrdSNEsHCg7HE4MLderL
noNj0zBlnpI5lbxMFPsFA2qhdGCGvpMiaOwbvsR9lz9QwcRYAASft9CCIp5LJc9t
qowrkn3DWFEkJhVkFTFJ8+Pvv5bMiAK1GFg1PhtgaK+t3ad7gDBf
=vGoy
-----END PGP PRIVATE KEY BLOCK-----
`

const publicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

mI0EVjQA3wEEAMIns0AwUleSC2Lq+nEJlKVmU7oiOqsSYcGfIuQijo5fW5ma0FJ/
XCVJUB69w0kkloeQtDPy8C8TuXPFW8cQvw/Z6DrNQaYlESUJt1w/oe2K9iC5q35V
1R7Tnsqnuw7KOGLs/ahRNSPsxrdEYC+iL3u6CAGXzyDBgvL6HPwwE96VABEBAAG0
CHRlc3Qga2V5iLgEEwECACIFAlY0AN8CGwMGCwkIBwMCBhUIAgkKCwQWAgMBAh4B
AheAAAoJEInlFUgfj6Eu5NUD/3EwqgEs+6d5T1mkhtvl0ECeeesMQIozEYX/1z6w
jJ9YXm9K6bIisHpUUwDqI3RuVsgFdpwGdxJOvmlbl4ieqkHwKzwV2I/IGJaEHpIZ
ldBqspLdxJkeT3ttM3egGgfruQ+UqcPgX4g2qyP56bSI6afdXC1ph6lMdJ3Pbh1x
IUkNuI0EVjQA3wEEALJBp2uuiQIgMO5UE2F51qkPBcOLUeY0W1wkj90oNGH/2POH
LxwEMFCQ1JQ3aqs6sf2K+DXMhCfCnbsetNWgPVIQQxSLqEEHqehE1J0c8KkVkJiv
4U/C+PcK5kapIysakThYCgytDtbx5wJYfZGXNdl44zw1ge5dhIlyuZICxqXrABEB
AAGInwQYAQIACQUCVjQA3wIbDAAKCRCJ5RVIH4+hLsXVBADApx070XiL1pDSss48
yuWtA2QpiS6BZifM5ja3UjRLBwoOxxODC3Xqy56DY9MwZZ6SOZW8TBT7BQNqoXRg
hr6TImjsG77EfZc/UMHEWAAEn7fQgiKeSyXPbaqMK5J9w1hRJCYVZBUxSfPj77+W
zIgCtRhYNT4bYGivrd2ne4AwXw==
=k72J
-----END PGP PUBLIC KEY BLOCK-----`