package main

import (
	"testing"
)

func TestShouldFindTotoInFile(m *testing.T) {
	filePath := "totofile.txt"
	pattern := "Toto"
	expected := []string{
		"Toto était un petit garçon curieux et intrépide.",
		"Chaque jour, Toto explorait les recoins les plus secrets de son quartier, à la recherche d'aventures excitantes.",
		"Un jour, Toto découvrit une vieille carte au trésor cachée sous une pile de feuilles mortes.",
		"Avec une énergie débordante, Toto se lança dans une quête incroyable pour trouver le trésor légendaire.",
		"À chaque étape de son périple, Toto rencontrait de nouveaux défis et de nouveaux amis.",
		"Finalement, après des jours d'exploration, Toto parvint à dénicher le trésor tant convoité, rempli de richesses inimaginables.",
		"Cette aventure extraordinaire marqua le début d'une série d'exploits épiques pour Toto, qui devint rapidement une légende dans sa petite ville.",
	}
	result, _ := grepFile(filePath, pattern, 4)

	if len(result) != len(expected) {
		m.Errorf("Expected %d results, got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			m.Errorf("Expected %s, got %s", expected[i], result[i])
		}
	}
}

func TestShouldFindJourInFile(m *testing.T) {
	filePath := "totofile.txt"
	pattern := "jour"
	expected := []string{
		"Chaque jour, Toto explorait les recoins les plus secrets de son quartier, à la recherche d'aventures excitantes.",
		"Un jour, Toto découvrit une vieille carte au trésor cachée sous une pile de feuilles mortes.",
		"Finalement, après des jours d'exploration, Toto parvint à dénicher le trésor tant convoité, rempli de richesses inimaginables.",
	}
	result, _ := grepFile(filePath, pattern, 4)

	if len(result) != len(expected) {
		m.Errorf("Expected %d results, got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			m.Errorf("Expected %s, got %s", expected[i], result[i])
		}
	}
}

func TestShouldBeCaseSensitive(m *testing.T) {
	filePath := "totofile.txt"
	pattern := "toto"
	expected := []string{}
	result, _ := grepFile(filePath, pattern, 4)

	if len(result) != len(expected) {
		m.Errorf("Expected %d results, got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			m.Errorf("Expected %s, got %s", expected[i], result[i])
		}
	}
}

func TestShouldFailOnNonExistingFile(m *testing.T) {
	filePath := "nonexistingfile.txt"
	pattern := "toto"
	expected := "open nonexistingfile.txt: The system cannot find the file specified."

	_, err := grepFile(filePath, pattern, 4)

	if err == nil {
		m.Error("Expected an error, got nil")
	}

	if err.Error() != expected {
		m.Errorf("Expected %s, got %s", expected, err.Error())
	}
}
