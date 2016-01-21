# dogsights
Create Datadog Metrics for Insights Queries

# Featrues:

* Go binary that runs off a config
  * ex:
  ```toml
    [general]
    api_key="YOUR_KEY_HERE"
    frequency="10s"
    
    [impressions]
    nrql="SELECT count(*) FROM AdServerEvents WHERE vungleType='reportAd'  since 2 hour ago  facet pub_app_id  LIMIT 20"
  ```
  
* auto builds in jenkins and creats a new pod with config changes
