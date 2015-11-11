require('./lib/ionic/js/ionic.bundle.min.js');
require('./lib/ionic-material/dist/ionic.material.min.js')
require('./lib/ng-websocket/ng-websocket.js');

angular.module('app', ['ionic', 'ngWebsocket', 'ionic-material'])

.run(function($ionicPlatform) {
  $ionicPlatform.ready(function() {
    // Hide the accessory bar by default (remove this to show the accessory bar above the keyboard
    // for form inputs)
    if(window.cordova && window.cordova.plugins.Keyboard) {
      cordova.plugins.Keyboard.hideKeyboardAccessoryBar(true);
      cordova.plugins.Keyboard.disableScroll(true);
    }
    if(window.StatusBar) {
      StatusBar.styleDefault();
    }
  });
})

.controller('MainCtrl', function($scope, $websocket, $http, $ionicLoading) {
  $scope.buttons = [0, 0, 0];
  $scope.host = document.location.host;
  $scope.connect = function connect(host) {
    $scope.host = host;
    $websocket
      .$new({
        url: "ws://" + host + "/ws",
        reconnect: false
      })
      .$on('$open', function onOpen() {
        $scope.$apply(function() { $scope.$emit('connected'); });
      })
      .$on('processing', function onProcessingClick(data) {
        $scope.$apply(function() { $scope.$emit('processing'); });
      })
      .$on('click', function onClick(data) {
        $scope.$apply(function() {
          $scope.buttons[data.button] = data.count;
          $scope.$emit('done');
        });
      })
      .$on('$close', function onClose() {
        $scope.$apply(function() {
          $scope.isConnected = false;
          $scope.$emit('done');
        });
      });
  }

  $scope.clickOn = function clickOn(index) {
    $scope.$emit('processing');
    $http.post('http://' + $scope.host + '/click/' + index);
  };

  $scope.$on('connected', function() {
    $http.get('http://' + $scope.host + '/clicks').then(function(response) {
      for (i = 0; i < response.data.length; i++) {
        $scope.buttons[i] = response.data[i];
      }
      $scope.isConnected = true;
    }, function() {
      alert("ERROR!");
    });
  });
  $scope.$on('processing', function(msg) {
    if ($scope.loading) return;
    $scope.loading = $ionicLoading.show({
      template: 'Processing click...',
      scope:    $scope,
    });
  });
  $scope.$on('done', function() {
    $ionicLoading.hide($scope.loading);
    delete($scope.loading);
  });
})
