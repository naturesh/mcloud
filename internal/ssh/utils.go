package ssh

import (
	"fmt"

	"github.com/naturesh/mcloud/internal/core"
)

func (c *Client) WaitForLog(pattern string, timeout int) error {
	cmd := fmt.Sprintf(`
		timeout %ds bash -c '
		until docker logs --tail 50 mcloud 2>&1 | grep -q "%s"; do
			sleep 2
		done'
	`, timeout, pattern)

	if err := c.Run(cmd); err != nil {
		return fmt.Errorf("%w: %s: %v", core.ErrNotFound, pattern, err)
	}

	return nil
}

func (c *Client) IsContainerRunning() bool {
	cmd := `[ "$(docker inspect -f '{{.State.Running}}' mcloud 2>/dev/null)" = "true" ]`
	err := c.Run(cmd)

	return err == nil
}
