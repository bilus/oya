Feature: Building

Background:
   Given I'm in project dir

# Scenario: No Oyafile
# Scenario: Missing job

Scenario: Successful build
  Given file ./Oyafile containing
    """
    jobs:
      all: |
        foo=4
        if [ $foo -ge 3 ]; then
          touch OK
        fi
        echo "Done"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Done

  """
  And file ./OK exists


Scenario: Nested Oyafiles
  Given file ./Oyafile containing
    """
    jobs:
      all: |
        touch Root
        echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    jobs:
      all: |
        touch Project1
        echo "Project1"
    """
  And file ./project2/Oyafile containing
    """
    jobs:
      all: |
        touch Project2
        echo "Project2"
    """
  When I run "oya build all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Root
  Project1
  Project2

  """
  And file ./Root exists
  And file ./Project1 exists
  And file ./Project2 exists

# Scenario: No rebuild
# Scenario: Minimal rebuild
# Scenario: Shell specification
