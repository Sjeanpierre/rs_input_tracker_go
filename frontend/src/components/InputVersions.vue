<style src="./styles/array_input.css"></style>
<script src="./script/input_versions.js"></script>

<template>
  <v-app id="inspire">
    <v-navigation-drawer width="210" fixed clipped class="blue darken-3 application theme--dark" app v-model="drawer">

      <!--start bullshit-->
      <v-list dense class="blue darken-3 info-list">
          <v-divider dark class="my-3"></v-divider>
          <v-list-tile>
            <v-list-tile-content>
              <v-list-tile-title class="white--text">
                {{this.account_name}}
              </v-list-tile-title>
              <v-icon color="white">arrow_downward</v-icon>
            </v-list-tile-content>
          </v-list-tile>
        <v-list-tile>
          <v-list-tile-content>
            <v-list-tile-title class="white--text">
              {{this.array_name}}
            </v-list-tile-title>
            <v-icon class="" color="white">arrow_downward</v-icon>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-content>
            <v-list-tile-title class="white--text">
              {{this.input_name}}
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
      <!--end bullshit-->


      <v-list dense class="blue darken-3">
        <!--left side menu-->
        <template v-for="(item, i) in menu_items">
          <v-layout row v-if="item.heading" align-center :key="i">
            <v-flex xs6>
              <v-subheader v-if="item.heading">
                {{ item.heading }}
              </v-subheader>
            </v-flex>
          </v-layout>
          <v-divider dark v-else-if="item.divider" class="my-3" :key="i"></v-divider>
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
      <v-spacer/>
    </v-toolbar>
    <v-content fluid fill-height justify-center align-center class="grey lighten-4">
      <v-card class="container">
        <v-card-title>
        <span class="headline blue-grey--text darken-3">
        Versions of {{this.input_name}}
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
          :loading="input_versions_loading"
          :rows-per-page-items="[10,20,30,{'text':'All','value':-1}]"
        >
          <template slot="items" slot-scope="props">
            <td class="text-xs-left">{{ props.item.version }}</td>
            <td class="text-xs-left">{{ props.item.type }}</td>
            <td class="text-xs-left">{{ props.item.created_at}}</td>
            <td class="text-xs-left">{{ props.item.value }}</td>
          </template>
          <v-alert slot="no-results" :value="true" color="error" icon="warning">
            Your search for "{{ search }}" found no results.
          </v-alert>
        </v-data-table>
      </v-card>
    </v-content>
  </v-app>
</template>
