#!/usr/bin/env bats

# !!! failing these tests mean work, because downstream 
# go fix the code & make test green OR adapt the test and release a major version
#

@test run_help {
  result=$(./semver -h)
  [ "$?" -eq 0  ]
}

@test run_via_pipe {
  result=$(echo v0 | ./semver)
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.0.0","major":"v0","majorminor":"v0.0","prerelease":"","build":"","source":"v0"}' ]
}

@test run_via_bash_redirection {
  result=$(./semver <<<"v0")
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.0.0","major":"v0","majorminor":"v0.0","prerelease":"","build":"","source":"v0"}' ]
}

@test run_via_file_redirection {
  result=$(./semver -release patch -format json < <(echo -e "v0.1.9\n"))
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.1.10","major":"v0","majorminor":"v0.1","prerelease":"","build":"","source":"v0.1.10"}' ]
}

@test run_minor {
  result=$(./semver <<<"v0.1")
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.1.0","major":"v0","majorminor":"v0.1","prerelease":"","build":"","source":"v0.1"}' ]
}

@test run_patch {
  result=$(./semver <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.1.9","major":"v0","majorminor":"v0.1","prerelease":"","build":"","source":"v0.1.9"}' ]
}

@test run_release_new_major {
  result=$(./semver -release major <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v1.0.0","major":"v1","majorminor":"v1.0","prerelease":"","build":"","source":"v1.0.0"}' ]
}

@test run_release_new_minor {
  result=$(./semver -release minor <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.2.0","major":"v0","majorminor":"v0.2","prerelease":"","build":"","source":"v0.2.0"}' ]
}

@test run_release_new_patch {
  result=$(./semver -release patch <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$result" = '{"canonical":"v0.1.10","major":"v0","majorminor":"v0.1","prerelease":"","build":"","source":"v0.1.10"}' ]
}

@test run_format_eval_major {
  eval $(./semver -release major -format eval <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$MAJOR" = 'v1' ]
  [ "$MAJORMINOR" = 'v1.0' ]
  [ "$CANONICAL" = 'v1.0.0' ]
}

@test run_format_eval_minor {
  eval $(./semver -release minor -format eval <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$MAJOR" = 'v0' ]
  [ "$MAJORMINOR" = 'v0.2' ]
  [ "$CANONICAL" = 'v0.2.0' ]
}

@test run_format_eval_patch {
  eval $(./semver -release patch-format eval <<<"v0.1.9")
  [ "$?" -eq 0  ]
  [ "$MAJOR" = 'v0' ]
  [ "$MAJORMINOR" = 'v0.1' ]
  [ "$CANONICAL" = 'v0.1.10' ]
}

