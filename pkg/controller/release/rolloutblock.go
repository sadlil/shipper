package release

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"

	shipper "github.com/bookingcom/shipper/pkg/apis/shipper/v1alpha1"
	rolloutBlockOverride "github.com/bookingcom/shipper/pkg/util/rolloutblock"
)

func (s *Scheduler) processRolloutBlocks(rel *shipper.Release) (shouldBlockRollout bool, nonOverriddenRBsStatement string) {
	relOverrideRB, ok := rel.Annotations[shipper.RolloutBlocksOverrideAnnotation]
	if !ok {
		relOverrideRB = ""
	}
	relOverrideRBs := rolloutBlockOverride.NewOverride(relOverrideRB)

	nsRBs, err := s.rolloutBlockLister.RolloutBlocks(rel.Namespace).List(labels.Everything())
	if err != nil {
		runtime.HandleError(fmt.Errorf("failed to list rollout block objects: %s", err))
	}

	gbRBs, err := s.rolloutBlockLister.RolloutBlocks(shipper.GlobalRolloutBlockNamespace).List(labels.Everything())
	if err != nil {
		runtime.HandleError(fmt.Errorf("failed to list rollout block objects: %s", err))
	}

	existingRBs := rolloutBlockOverride.NewOverrideFromRolloutBlocks(append(nsRBs, gbRBs...))
	nonExistingRbs := relOverrideRBs.Diff(existingRBs)
	if len(nonExistingRbs) > 0 {
		for o := range nonExistingRbs {
			s.removeRolloutBlockFromAnnotations(relOverrideRBs, o, rel)
		}
		s.recorder.Event(rel, corev1.EventTypeWarning, "Non Existing RolloutBlock", nonExistingRbs.String())
	}

	nonOverriddenRBs := existingRBs.Diff(relOverrideRBs)
	shouldBlockRollout = len(nonOverriddenRBs) != 0
	nonOverriddenRBsStatement = nonOverriddenRBs.String()

	if shouldBlockRollout {
		s.recorder.Event(rel, corev1.EventTypeNormal, "RolloutBlock", nonOverriddenRBsStatement)
	} else if len(relOverrideRB) > 0 {
		s.recorder.Event(rel, corev1.EventTypeNormal, "Overriding RolloutBlock", relOverrideRB)
	}

	return
}

func (s *Scheduler) removeRolloutBlockFromAnnotations(overrideRBs rolloutBlockOverride.Override, rbName string, release *shipper.Release) {
	overrideRBs.Delete(rbName)
	release.Annotations[shipper.RolloutBlocksOverrideAnnotation] = overrideRBs.String()
	_, err := s.clientset.ShipperV1alpha1().Releases(release.Namespace).Update(release)
	if err != nil {
		runtime.HandleError(err)
	}
}