This program calculates the apdex from Dynatrace's user session export feed.

This is necessary because Dynatrace does not calculate the apdex for native
mobile applications yet (as of May 2018).

It feeds the apdex back using Dynatrace's custom network device API.


  ---------   user session export
 |Dynatrace|----------------------
  ---------                       |
      ^                           |
      |                           V
      |                       ----------
      |                      |  Apdex   |
      |                      |Calculator|
      |                       ----------
      |                           |
      |                           |
       ---------------------------
        custom network device API

To use this,
1. Clone this with git clone - this project uses Zeal Jagannatha's golang-ring
package. Don't forget to use the --recurse-submodules option when running git
clone.
2. Generate an API token in Dynatrace: Settings / Integration / Dynatrace API
3. Go to Settings / Integration / User session export.
4. Enable user session export.
5. Set the endpoint URL to the apdex calculator.
6. Build a docker image by running the buildimage script.
7. Edit the run script and modify the SERVERURL and APITOKEN variables.
8. Run the container with the run script.

