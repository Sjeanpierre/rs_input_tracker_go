<style src="./styles/array_input.css"></style>
<script src="./script/array_inputs.js"></script>

<template>
  <v-app id="inspire">
    <v-navigation-drawer width="210" fixed clipped class="blue darken-3 application theme--dark" app v-model="drawer">
      <v-list dense class="blue darken-3">
        <!--account info-->
        <v-divider dark class="my-3"></v-divider>
        <v-list-tile>
          <div style="display: block;margin:  auto;">
            <div class="info-list subheading">{{this.account_name}}</div>
            <div class="info-list"><v-flex offset>Array Count: {{this.array_count}}</v-flex></div>
            <div class="info-list offset-xs1">
              <v-tooltip bottom light>
                <v-flex offset slot="activator">Refreshed: {{array_last_updated_date}}</v-flex>
                <span>{{array_last_updated_date_long}}</span>
              </v-tooltip>
            </div>
          </div>
        </v-list-tile>
      </v-list>
        <v-list dense class="blue darken-3">
        <template v-for="(item, i) in menu_items">
          <v-layout row v-if="item.heading" align-center :key="i">
            <v-flex xs6>
              <v-subheader v-if="item.heading">
                {{ item.heading }}
              </v-subheader>
            </v-flex>
          </v-layout>
          <v-list-tile :key="i" v-else-if="item.snackbarToogle" @click="" @click.native="snackbar = !snackbar">
            <v-list-tile-action>
              <v-icon color="white">archive</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title class="white--text">
                Force Refresh
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-divider dark v-else-if="item.divider" class="my-3" :key="i"/>
          <v-list-tile :key="i" v-else @click="" :to="item.link">
            <v-list-tile-action>
              <v-icon color="white">{{ item.icon }}</v-icon>
            </v-list-tile-action>
            <v-list-tile-content>
              <v-list-tile-title class="white--text">
                {{ item.text }}
              </v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </template>
        <!--end left side menu-->
      </v-list>
    </v-navigation-drawer>
    <v-toolbar height="35px" color="white" app absolute clipped-left>
      <v-toolbar-side-icon @click.native="drawer = !drawer"/>
      <site-title/>
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-content fluid fill-height justify-center align-center class="grey lighten-4">
              <v-card class="container">
      <v-card-title>
        <span class="headline blue-grey--text darken-3">
        {{this.array_name}}
        </span>
        <v-spacer/>
        <v-text-field
          append-icon="search"
          label="Search Inputs"
          single-line
          hide-details
          v-model="search"
        />
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="items"
        :search="search"
        :loading="datatable_loading"
        :rows-per-page-items = "[10,20,30,{'text':'All','value':-1}]"
      >
        <template slot="items" slot-scope="props">
          <td class="text-xs-left"><v-btn v-bind="inputVersionLink(props.item.name)" color="primary" flat>{{ props.item.name }}</v-btn></td>
          <td class="text-xs">{{ props.item.type }}</td>
          <td class="text-xs">{{ props.item.version }}</td>
          <td class="text-xs-left">{{ props.item.value }}</td>
        </template>
        <v-alert slot="no-results" :value="true" color="error" icon="warning">
          Your search for "{{ search }}" found no results.
        </v-alert>
      </v-data-table>
    </v-card>
    </v-content>
    <v-snackbar :timeout="6000" color='error' class="title" v-model="snackbar">Sorry, this feature is not yet implemented 🙈
      <v-btn dark flat @click.native="snackbar = false">Close</v-btn></v-snackbar>
  </v-app>
</template>
