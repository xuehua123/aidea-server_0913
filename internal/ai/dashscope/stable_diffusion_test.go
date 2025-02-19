package dashscope_test

import (
	"context"
	"testing"
	"time"

	"github.com/mylxsw/aidea-server/internal/ai/dashscope"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-utils/assert"
)

func TestStableDiffusion(t *testing.T) {
	client := createClient()

	resp, err := client.StableDiffusion(context.TODO(), dashscope.StableDiffusionRequest{
		Model: dashscope.ImageModelSDXL,
		Input: dashscope.StableDiffusionInput{
			Prompt: "no humans, outdoors, scenery, power lines, clouds, sky, poles, 1 houses, buildings, grass, plants, blue sky, cloudy sky, windows, trees, day, distant view, nice clouds, depth of field, distant view, (Illustration: 1.0), masterpiece, best quality",
		},
		Parameters: dashscope.StableDiffusionParameters{
			N: 1,
		},
	})
	assert.NoError(t, err)

	log.With(resp).Debug("resp")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for {
		select {
		case <-ticker.C:
			taskResp, err := client.ImageTaskStatus(context.TODO(), resp.Output.TaskID)
			assert.NoError(t, err)

			log.With(taskResp).Debug("resp")

			if taskResp.Output.TaskStatus == dashscope.TaskStatusSucceeded || taskResp.Output.TaskStatus == dashscope.TaskStatusFailed {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
