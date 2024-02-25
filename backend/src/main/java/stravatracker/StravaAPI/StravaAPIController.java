package stravatracker.StravaAPI;


import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

@RestController
public class StravaAPIController extends RestTemplate {

    StravaAPIService stravaAPIService;

    @Autowired
    public StravaAPIController(StravaAPIService stravaAPIService) {
        this.stravaAPIService = stravaAPIService;
    }

    public JsonNode getClubActivities() {
        return stravaAPIService.getClubActivities();
    }

    public JsonNode getClubMembers() {
        return stravaAPIService.getClubMembers();
    }


}
