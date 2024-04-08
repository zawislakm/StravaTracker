package stravatracker.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import jakarta.persistence.*;
import lombok.*;

import java.util.ArrayList;
import java.util.List;

@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
@Entity(name = "athletes")
@Table(name = "athletes",
        uniqueConstraints = {
                @UniqueConstraint(name = "athlete_firstname_lastname_unique", columnNames = {"firstname", "lastname"})
        }
)
@ToString(exclude = "activities")
@EqualsAndHashCode(exclude = "activities")
public class Athlete {

    @Id
    @SequenceGenerator(name = "athlete_sequence", sequenceName = "athlete_sequence", allocationSize = 1)
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "athlete_sequence")
    @Column(name = "id", updatable = false, unique = true, nullable = false)
    @JsonIgnore
    private long id;

    @Column(name = "firstname", columnDefinition = "TEXT", nullable = false)
    private String firstName;

    @Column(name = "lastname", columnDefinition = "TEXT", nullable = false)
    private String lastName;

    @OneToMany(mappedBy = "athlete", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    @JsonIgnoreProperties("activities")
    @JsonIgnore
    private final List<Activity> activities = new ArrayList<>();

    public Athlete(String firstName, String lastName) {
        this.firstName = firstName;
        this.lastName = lastName;
    }

}

