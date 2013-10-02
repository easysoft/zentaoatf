<?php
/**
 * This file contains functions to ouput case info in yaml format.
 */

/**
 * Print case title.
 * 
 * @param  string    $title 
 * @access public
 * @return void
 */
function title($title)
{
    echo "title: $title\n";
}

/**
 * Run a statement and print the output.
 * 
 * @param  string    $output 
 * @access public
 * @return void
 */
function run($output)
{
    static $firstCall = true;

    if($firstCall) (print "steps: \n") && ($firstCall = false);

    echo space(2) . "-\n";

    $output  = trim($output);
    $lines = explode("\n", $output);

    if(count($lines) == 1) return print(space(4) . "output: $output\n");

    echo space(4) . "output: |\n";
    foreach($lines as $line) echo space(6) . trim($line) . "\n";
}

/**
 * Print the expect section.
 * 
 * @param  string    $expect 
 * @param  string    $tag        expect|expectx
 * @access public
 * @return void
 */
function expect($expect, $tag = 'expect')
{
    $expect  = trim($expect);
    $lines = explode("\n", $expect);

    if(count($lines) == 1) return print(space(4) . "$tag: $expect\n");

    echo space(4) . "$tag: |\n";
    foreach($lines as $line) echo space(6) . trim($line) . "\n";
}

/**
 * Print expectx section.
 * 
 * @param  int    $expect 
 * @access public
 * @return void
 */
function expectx($expect)
{
    expect($expect, 'expectx');
}

/**
 * Print some spaces.
 * 
 * @param  int    $number 
 * @access public
 * @return string
 */
function space($number)
{
    return str_repeat(' ', $number);
}
