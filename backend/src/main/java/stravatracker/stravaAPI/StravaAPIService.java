package stravatracker.stravaAPI;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;
import stravatracker.service.ActivityService;
import stravatracker.service.AthleteService;
import stravatracker.service.SportTypeService;
import stravatracker.model.Athlete;
import stravatracker.model.SportType;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;

@Service
public class StravaAPIService {

    @Value("${strava.CLUB_ID}")
    private int CLUB_ID = 1131273;
    WebClient webClient;
    ActivityService activityService;
    AthleteService athleteService;
    SportTypeService sportTypeService;

    @Autowired
    public StravaAPIService(WebClient webClient, ActivityService activityService, AthleteService athleteService, SportTypeService sportTypeService) {
        this.webClient = webClient;
        this.activityService = activityService;
        this.athleteService = athleteService;
        this.sportTypeService = sportTypeService;
    }


    @Scheduled(fixedDelay = 60 * 60 * 1000) // once per hour
    public void processClubActivities() {
        JsonNode newActivities = this.getClubActivities();
        if (newActivities == null) return;

        JsonNode lastActivity = this.getLastSavedClubActivity();
        for (JsonNode activity : newActivities) {

            if (activity.equals(lastActivity)) {
                break;
            }
            Athlete athlete = athleteService.getOrCreateAthlete(activity.path("athlete"));
            SportType sportType = sportTypeService.getOrCreateSportType(activity);
            activityService.addNewActivity(activity, athlete, sportType);
        }

        this.saveLastSavedClubActivity(newActivities.get(0));
    }


    @Scheduled(fixedDelay = 24 * 60 * 60 * 1000) // once per day
    public void processClubMembers() {
        JsonNode clubMembers = this.getClubMembers();
        if (clubMembers == null) return;

        for (JsonNode member : clubMembers) {
            athleteService.getOrCreateAthlete(member);
        }
    }


    public void saveLastSavedClubActivity(JsonNode lastActivity) {
        String filePath = "src/main/resources/lastSavedClubActivity.json";
        try {
            ObjectMapper mapper = new ObjectMapper();
            mapper.writeValue(new File(filePath), lastActivity);
        } catch (IOException e) {
            throw new RuntimeException("Error writing JSON file: " + filePath + e.getMessage(), e);
        }

    }

    public JsonNode getLastSavedClubActivity() {
        // Strava doesn't provide any id of activity (in club api responses),\
        // so to remember last activity is needed to save it in file (to avoid duplicates in the database)
        String filePath = "src/main/resources/lastSavedClubActivity.json";
        try {
            ObjectMapper mapper = new ObjectMapper();
            return mapper.readTree(new File(filePath));
        } catch (FileNotFoundException e) {
            System.out.println("File not found: " + filePath);
            return null;
        } catch (IOException e) {
            throw new RuntimeException("Error reading JSON file: " + filePath + e.getMessage(), e);

        }
    }

    public JsonNode getClubActivities() {
        return webClient.get()
                .uri("https://www.strava.com/api/v3/clubs/" + CLUB_ID + "/activities")
                .retrieve()
                .bodyToMono(JsonNode.class)
                .doOnError(throwable -> System.err.println("Error occurred while fetching club members from Strava API: " + throwable.getMessage()))
                .onErrorResume(throwable -> Mono.empty())
                .block();
    }


    public JsonNode getClubMembers() {
        return webClient.get()
                .uri("https://www.strava.com/api/v3/clubs/" + CLUB_ID + "/members")
                .retrieve()
                .bodyToMono(JsonNode.class)
                .doOnError(throwable -> System.err.println("Error occurred while fetching club members from Strava API: " + throwable.getMessage()))
                .onErrorResume(throwable -> Mono.empty())
                .block();
    }
}
