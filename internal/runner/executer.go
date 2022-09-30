package runner

import (
	"github.com/pkg/errors"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/katana/pkg/engine"
	"github.com/projectdiscovery/katana/pkg/engine/hybrid"
	"github.com/projectdiscovery/katana/pkg/engine/standard"
	"github.com/remeh/sizedwaitgroup"
)

// ExecuteCrawling executes the crawling main loop
func (r *Runner) ExecuteCrawling() error {
	inputs := r.parseInputs()
	if len(inputs) == 0 {
		return errors.New("no input provided for crawling")
	}

	var (
		crawler engine.Engine
		err     error
	)

	switch {
	case r.options.Headless:
		crawler, err = hybrid.New(r.crawlerOptions)
	default:
		crawler, err = standard.New(r.crawlerOptions)
	}
	if err != nil {
		return errors.Wrap(err, "could not create standard crawler")
	}
	defer crawler.Close()

	wg := sizedwaitgroup.New(r.options.Parallelism)
	for _, input := range inputs {
		wg.Add()

		go func(input string) {
			defer wg.Done()

			if err := crawler.Crawl(input); err != nil {
				gologger.Warning().Msgf("%s\n", err)
			}
		}(input)
	}
	wg.Wait()
	return nil
}
