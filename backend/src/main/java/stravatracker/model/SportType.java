package stravatracker.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonManagedReference;
import jakarta.persistence.*;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.Setter;
import lombok.ToString;

import java.util.ArrayList;
import java.util.List;


@Entity(name = "sport_types")
@Table(name = "sport_types", uniqueConstraints = {
        @UniqueConstraint(name = "sport_type_unique", columnNames = {"type", "sport_type"})
})

@ToString(exclude = "activities")
@EqualsAndHashCode(exclude = "activities")
public class SportType {

    @Id
    @SequenceGenerator(name = "sport_type_sequence", sequenceName = "sport_type_sequence", allocationSize = 1)
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "sport_type_sequence")
    @Column(name = "id", updatable = false, unique = true, nullable = false)
    @JsonIgnore
    private Long id;


    @Column(name = "type", columnDefinition = "TEXT", nullable = false)
    private String type;

    @Column(name = "sport_type", columnDefinition = "TEXT", nullable = false)
    private String sportType;

    @OneToMany(mappedBy = "sportType", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
//    @JsonManagedReference
    @JsonIgnore
    private final List<Activity> activities = new ArrayList<>();

    public SportType() {

    }

    public SportType(String type, String sportType) {
        this.type = type;
        this.sportType = sportType;
    }

    public SportType(Long id, String type, String sportType) {
        this.id = id;
        this.type = type;
        this.sportType = sportType;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getSportType() {
        return sportType;
    }

    public void setSportType(String sportType) {
        this.sportType = sportType;
    }

    public List<Activity> getActivities() {
        return activities;
    }
}
