Feature: Running tasks

Background:
   Given I'm in project dir

Scenario: Import tasks from vendored packs
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    all: |
      foo=4
      if [ $$foo -ge 3 ]; then
        touch OK
      fi
      echo "Done"
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Done

  """
  And file ./OK exists

Scenario: Import task using pack values
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    Values:
      foo: xxx
    all: |
      bar=$foo
      echo $$bar
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  xxx

  """

Scenario: Import task using BasePath
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    Values:
      foo: xxx
    all: |
      bar=$$(basename $BasePath)
      echo $$bar
    """
  When I run "oya run foo.all"
  Then the command succeeds
  And the command outputs to stdout
  """
  foo

  """

Scenario: Access pack values
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    all: |
      echo $foo.bar
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    Values:
      bar: xxx
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  xxx

  """

Scenario: Access current project values
  Given file ./Oyafile containing
    """
    Project: main
    Values:
      foo: main
    """
  And file ./project1/Oyafile containing
    """
    Values:
      foo: project1
    """
  And file ./project2/Oyafile containing
    """
    Import:
      main: /
      p1: /project1
    Values:
      foo: project2
    all: |
      echo $main.foo
      echo $p1.foo
      echo $foo
    """
  When I run "oya run all"
  Then the command succeeds
  And the command outputs to stdout
  """
  main
  project1
  project2

  """

Scenario: Invalid import
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo

    all: echo "OK"
    """
  When I run "oya run all"
  Then the command fails with error matching
  """
  .*missing pack github.com/test/foo$
  """