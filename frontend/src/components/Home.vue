<style src="./styles/home.css"></style>
<script src="./script/home.js"></script>


<template>
  <v-app id="inspire">
    <!--<v-toolbar-title class="display-3 text-center mt-3">Rightscale Input Auditor</v-toolbar-title>-->
    <v-layout row mt-5>
      <v-flex xs12 sm4 offset-sm4>
        <v-card>
          <v-toolbar height="40px" color="primary" dark>
            <v-toolbar-title class="text-center">Rightscale Accounts 😐</v-toolbar-title>
            <v-spacer></v-spacer>
          </v-toolbar>
          <v-list subheader>
          <v-list-tile avatar v-for="item in items" :key="item.name" @click="">
            <v-list-tile-avatar color="primary">
              <v-icon dark>cloud_queue</v-icon>
            </v-list-tile-avatar>
            <v-list-tile-content>
              <v-list-tile-title class="subheading account-name" v-html="item.name"></v-list-tile-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-btn v-bind="arrayLink(item.account_id)" outline color="indigo">Array List</v-btn>
            </v-list-tile-action>
          </v-list-tile>
          </v-list>
        </v-card>
        <v-toolbar height="40px" color="primary" dark>
          <v-toolbar-title class="text-center"></v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>
        <!--<v-footer height="auto">-->
          <!--<v-card flat tile class="flex">-->
            <!--<v-card-title class="white&#45;&#45;text primary">-->
              <!--&lt;!&ndash;<strong class="subheading">Get connected with us on social networks!</strong>&ndash;&gt;-->
            <!--</v-card-title>-->
          <!--</v-card>-->
        <!--</v-footer>-->
      </v-flex>
    </v-layout>
  </v-app>
</template>
