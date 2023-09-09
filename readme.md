# Intro

Projects feed is a project aggregator and RSS provider from various (freelancer) websites. Available sources include:

- [Project.co.id](https://projects.co.id)
- [Sribu.com](https://www.sribu.com)

### Response Content

There are 3 different available content, HTML, JSON, and Feed (RSS feed, ATOM feed, & JSON feed).

The HTML response is provided for end users to view available projects directly from their browsers. By accessing the domain [projects-feed.fly.io](https://projects-feed.fly.io), you can see the default project as HTML.

RSS provides formatted XML for feeds. By accessing URL [projects-feed.fly.io/projects/rss](https://projects-feed.fly.io/projects/rss), you can see XML that you can use to subscribe to the feed. Different formats for the feed are available. You can change the format by changing the suffix based on the type. For example:

```bash
https://projects-feed.fly.io/projects/atom: will respond as an RSS feed
https://projects-feed.fly.io/projects/json: will respond as a JSON feed
```
While JSON feed has limited data on structure, you can also fetch projects as complete JSON by accessing the URL [projects-feed.fly.io/projects](https://projects-feed.fly.io/projects). The response provides complete enough data if you want to develop different features based on the feed.

> **Note:**
> Every content data is cached for 15 minutes to reduce requests to project providers and improve performance.

### Filtering

By default, each content response will return a list of projects from all vendors. It includes all projects from all categories sorted by the latest published date.

If you want to filter the response for a specific vendor, you can add the query `vendor=name` to the URL. For example:

- For HTML: [https://projects-feed.fly.dev/?vendor=Projects.co.id](https://projects-feed.fly.dev/?vendor=Projects.co.id)
- For JSON response (API): [https://projects-feed.fly.dev/projects/?vendor=Sribu.com](https://projects-feed.fly.dev/projects/?vendor=Sribu.com)
- For feed: [https://projects-feed.fly.dev/projects/rss?vendor=Projects.co.id](https://projects-feed.fly.dev/projects/rss?vendor=Projects.co.id)

You can also filter projects by their category by adding the query `tag=name` to the URL. Here are example generated URLs to filter by category:

- For HTML: [https://projects-feed.fly.io?tag=SEO](https://projects-feed.fly.io?tag=SEO)
- For JSON response (API): [https://projects-feed.fly.dev/projects?tag=Video%20Editing](https://projects-feed.fly.dev/projects?tag=Video%20Editing)
- For feed: [https://projects-feed.fly.dev/projects/rss?tag=Video%20Editing](https://projects-feed.fly.dev/projects/rss?tag=Video%20Editing)

You can combine both filters based on your needs. For example:

- To only show projects with the "data entry" tag from "projects.co.id": [https://projects-feed.fly.dev/?tag=Data%20Entry&vendor=Projects.co.id](https://projects-feed.fly.dev/?tag=Data%20Entry&vendor=Projects.co.id)
- To only show the feed for marketing from "projects.co.id" in RSS feed: [https://projects-feed.fly.dev/projects/rss?tag=Data%20Entry&vendor=projects.co.id](https://projects-feed.fly.dev/projects/rss?tag=Data%20Entry&vendor=projects.co.id)

These filters are applied to all kinds of responses, including JSON, HTML, and feed.