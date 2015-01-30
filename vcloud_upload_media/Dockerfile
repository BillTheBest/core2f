FROM ruby:2.1
MAINTAINER EMCCODE

RUN gem install nokogiri netaddr
RUN gem install httpclient rest-client
RUN gem install builder

WORKDIR /
RUN git clone https://github.com/emccode/ruby_vcloud_sdk
WORKDIR /ruby_vcloud_sdk
RUN git checkout updates-vcaondemand-uploadmedia && gem build ruby_vcloud_sdk.gemspec && gem install --local ruby_vcloud_sdk

ADD ./vcloud_upload_media.rb /vcloud_upload_media.rb

CMD ruby /vcloud_upload_media.rb

