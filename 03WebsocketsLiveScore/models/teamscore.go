/*
* Created on 02 March 2024
* @author Sai Sumanth
 */
package models

type TeamScore struct {
	TeamName     string `json:"team_name"`
	TotalScore   int    `json:"total_score"`
	TotalWickets int    `json:"total_wickets"`
}