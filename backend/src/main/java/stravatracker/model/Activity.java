package stravatracker.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.databind.JsonNode;
import jakarta.persistence.*;
import lombok.*;

import java.time.LocalDate;

@NoArgsConstructor
@AllArgsConstructor
@Setter
@Getter
@Entity(name = "activities")
@ToString(exclude = {"id", "athlete", "sportType"})
@EqualsAndHashCode(exclude = {"id", "athlete", "sportType"})
public class Activity {

    @Id
    @SequenceGenerator(name = "activity_sequence", sequenceName = "activity_sequence", allocationSize = 1)
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "activity_sequence")
    @Column(name = "id", updatable = false, unique = true, nullable = false)
    @JsonIgnore
    private Long id;

    @ManyToOne
    @JoinColumn(name = "athlete_id", nullable = false)
    private Athlete athlete;

    @ManyToOne
    @JoinColumn(name = "sport_type_id", nullable = false)
    private SportType sportType;

    @Column(name = "distance", nullable = false)
    private Float distance;

    @Column(name = "moving_time", nullable = false)
    private Float movingTime;

    @Column(name = "total_elevation_gain", nullable = false)
    private Float totalElevationGain;

    @Column(name = "activity_date", nullable = false, columnDefinition = "DATE")
    private LocalDate activityDate;

    public Activity(Athlete athlete, SportType sportType, Float distance, Float movingTime, Float elevationGain, LocalDate activityDate) {
        this.athlete = athlete;
        this.sportType = sportType;
        this.distance = distance;
        this.movingTime = movingTime;
        this.totalElevationGain = elevationGain;
        this.activityDate = activityDate;
    }

    public Activity(SportType sportType, Athlete athlete) {
        this.sportType = sportType;
        this.athlete = athlete;
        this.distance = 0f;
        this.movingTime = 0f;
        this.totalElevationGain = 0f;
    }

    public Activity(JsonNode jsonData, Athlete athlete, SportType sportType) {
        this.athlete = athlete;
        this.sportType = sportType;
        this.distance = jsonData.get("distance").floatValue();
        this.movingTime = jsonData.get("moving_time").floatValue();
        this.totalElevationGain = jsonData.get("total_elevation_gain").floatValue();
        this.activityDate = LocalDate.now();
    }

    public Float getDistance() {
        return (distance == null) ? 0f : this.distance;
    }

    public Float getMovingTime() {
        return (movingTime == null) ? 0f : this.movingTime;
    }

    public Float getTotalElevationGain() {
        return (totalElevationGain == null) ? 0f : this.totalElevationGain;
    }

}
