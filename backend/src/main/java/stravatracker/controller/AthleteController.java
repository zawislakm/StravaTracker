package stravatracker.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import stravatracker.service.AthleteService;
import stravatracker.model.Athlete;

import java.util.List;

@RestController
@RequestMapping("/api/v1/athlete")
public class AthleteController {


    private final AthleteService athleteService;

    @Autowired
    public AthleteController(AthleteService athleteService) {
        this.athleteService = athleteService;
    }

    @GetMapping("/athletes")
    public List<Athlete> getAllAthletes() {
//        this.athleteService.addNewAthlete(new Athlete("Maks","Z"));
        return this.athleteService.getAthletes();

    }

//    @GetMapping("/activities")
//    public List<>
}
