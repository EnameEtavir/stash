package autotag

import (
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getMatchingPerformers(path string, performerReader models.PerformerReader) ([]*models.Performer, error) {
	words := getPathWords(path)
	performers, err := performerReader.QueryForAutoTag(words)

	c := config.GetInstance()

	if err != nil {
		return nil, err
	}

	var ret []*models.Performer

	// TODO: currently we only conservatively only match "qualified" names which are judged to be good for matching.
	// 		As a follow-up there should be "possible" matches added (unique match via "qualified" and "possible") that
	//      also delivers the UI feature to move a match between "possible" and "qualified".
	// 		For example a second getPossiblyMatchingPerformers for performers_possible_scenes,
	//		performers_possible_images ... should be added
	for _, p := range performers {
		qualifiedNames := getQualifiedPerformerNames(p)
		for _, qualifiedName := range qualifiedNames {
			if qualifiedName.Qualified || !c.GetOnlyQualifiedPerformers() {
				if nameMatchesPath(qualifiedName.String, path) {
					ret = append(ret, p)
				}
			}
		}
	}

	return ret, nil
}

func getPerformerTagger(p *models.Performer) tagger {
	return tagger{
		ID:      p.ID,
		Type:    "performer",
		Name:    getQualifiedPerformerName(p.Name.String),
		Aliases: getQualifiedPerformerAliases(p),
	}
}

// PerformerScenes searches for scenes whose path matches the provided performer name and tags the scene with the performer.
func PerformerScenes(p *models.Performer, paths []string, rw models.SceneReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(rw, otherID, subjectID)
	})
}

// PerformerImages searches for images whose path matches the provided performer name and tags the image with the performer.
func PerformerImages(p *models.Performer, paths []string, rw models.ImageReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagImages(paths, rw, func(subjectID, otherID int) (bool, error) {
		return image.AddPerformer(rw, otherID, subjectID)
	})
}

// PerformerGalleries searches for galleries whose path matches the provided performer name and tags the gallery with the performer.
func PerformerGalleries(p *models.Performer, paths []string, rw models.GalleryReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagGalleries(paths, rw, func(subjectID, otherID int) (bool, error) {
		return gallery.AddPerformer(rw, otherID, subjectID)
	})
}
