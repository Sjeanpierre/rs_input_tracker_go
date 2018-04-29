<style src="./styles/array_input.css"></style>
<script src="./script/array_list.js"></script>

<template>
  <v-app id="inspire">
    <v-toolbar height="35px" color="white" app absolute clipped-left>
      <site-title/>
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-content fluid fill-height justify-center align-center class="grey lighten-4">
      <v-card width="60%" class="container">
        <v-card-title>
        <span class="headline blue-grey--text darken-3">
        {{this.account_name}} Arrays
        </span>
          <v-spacer></v-spacer>
          <v-text-field
            append-icon="search"
            label="Search Array"
            single-line
            hide-details
            v-model="search"
          ></v-text-field>
        </v-card-title>
        <v-data-table
          :headers="headers"
          :items="items"
          :loading="array_list_loading"
          :search="search"
          :rows-per-page-items="[10,20,30,{'text':'All','value':-1}]"
        >
          <template slot="items" slot-scope="props">
            <td class="text-xs-left">{{ props.item.array_name }}</td>
            <td class="text-xs-center">
              <div class="text-xs-center">
                <v-btn
                  :to="{ name: 'ArrayInputs', params: { account_id: props.item.account_id, array_id: props.item.array_id }}"
                  outline color="green">{{ props.item.array_id }}
                </v-btn>
              </div>
            </td>
            <td class="text-xs-center">
              <v-flex>
                <v-btn flat icon color="blue" v-bind="rightscaleLink(props.item.array_id)">
                  <v-icon>link</v-icon>
                </v-btn>
              </v-flex>
            </td>
          </template>
          <v-alert slot="no-results" :value="true" color="error" icon="warning">
            Your search for "{{ search }}" found no results.
          </v-alert>
        </v-data-table>
      </v-card>
    </v-content>
  </v-app>
</template>
