package stravatracker.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import stravatracker.service.ActivityService;
import stravatracker.model.Activity;
import stravatracker.model.Athlete;
import stravatracker.requestsmodels.ActivitySumUpRequest;
import stravatracker.requestsmodels.ActivitySumUpResponse;
import stravatracker.requestsmodels.DateRange;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/activity")
public class ActivityController {

    private final ActivityService activityService;

    @Autowired
    public ActivityController(ActivityService activityService) {
        this.activityService = activityService;
    }

    @GetMapping("/activities")
    public List<Activity> getActivities() {
        return this.activityService.getActivities();
    }

    @PostMapping("/activities")
    public List<Activity> getActivitiesBetweenDates(@RequestBody DateRange dateRange) {
        return this.activityService.getActivities(dateRange);
    }


    @PostMapping(value = "/athletes/sumUp")
    @ResponseBody
    public Map<Athlete, ActivitySumUpResponse> getActivitiesSumUpByAthletes(@RequestBody ActivitySumUpRequest activitySumUpRequest){
        System.out.println("-------------------------------------------------------------");
        System.out.println(activitySumUpRequest.toString());
        return this.activityService.getActivitiesSumUpByAthletes(activitySumUpRequest.getSportType(), activitySumUpRequest.getDateRange());
    }
}
